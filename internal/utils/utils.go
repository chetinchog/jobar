package utils

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

// --- Regex compiladas para limpieza HTML ---
var (
	// Bloques completos que se eliminan enteros (script, style, etc.)
	reScript   = regexp.MustCompile(`(?is)<script[\s\S]*?</script>`)
	reStyle    = regexp.MustCompile(`(?is)<style[\s\S]*?</style>`)
	reComment  = regexp.MustCompile(`(?s)<!--[\s\S]*?-->`)

	// Tags de bloque que se convierten a salto de línea antes de quitar el HTML.
	// Incluye: p, br, div, li, h1-h6, tr, hr (con o sin atributos, apertura y cierre).
	reBlockTag = regexp.MustCompile(`(?i)</?\s*(?:p|br|div|li|ul|ol|h[1-6]|tr|td|th|hr|blockquote|pre|section|article|header|footer|aside|nav)[^>]*>`)

	// Cualquier tag HTML restante con o sin atributos (incluyendo URLs dentro)
	reHTMLTag  = regexp.MustCompile(`<[^>]+>`)

	// URLs sueltas en el texto (http/https/ftp)
	reURL      = regexp.MustCompile(`https?://\S+|ftp://\S+`)

	// HTML entities numéricas como &#160; o &#x00A0;
	reEntityNum = regexp.MustCompile(`&#x?[0-9a-fA-F]+;`)

	// Espacios/tabs múltiples en una misma línea
	reMultiSpace = regexp.MustCompile(`[ \t]{2,}`)

	// Más de 2 saltos de línea consecutivos → los aplana a 2
	reMultiNewline = regexp.MustCompile(`\n{3,}`)
)

// htmlEntities mapea las entidades HTML más comunes a sus equivalentes en texto plano.
var htmlEntities = strings.NewReplacer(
	"&nbsp;",  " ",
	"&amp;",   "&",
	"&lt;",    "<",
	"&gt;",    ">",
	"&quot;",  `"`,
	"&apos;",  "'",
	"&#39;",   "'",
	"&ndash;", "-",
	"&mdash;", "-",
	"&ldquo;", `"`,
	"&rdquo;", `"`,
	"&lsquo;", "'",
	"&rsquo;", "'",
	"&bull;",  "-",
	"&middot;","·",
	"&hellip;","...",
	"&copy;",  "(c)",
	"&reg;",   "(R)",
	"&trade;", "(TM)",
	"&euro;",  "€",
	"&pound;", "£",
	"&yen;",   "¥",
	"&cent;",  "¢",
	"&laquo;", "«",
	"&raquo;", "»",
	"&frac12;","1/2",
	"&frac14;","1/4",
	"&frac34;","3/4",
)

// CleanForSearch limpia texto para uso en búsquedas (sin HTML, todo minúsculas).
func CleanForSearch(s string) string {
	s = CleanHTML(s)
	return strings.ToLower(s)
}

// CleanHTML elimina todo el HTML (tags, scripts, styles, comentarios, URLs,
// entities) y devuelve texto plano limpio preservando saltos de línea
// en los tags de bloque (p, br, li, div, h1-h6, etc.).
//
// IMPORTANTE: primero decodifica las HTML entities para manejar correctamente
// los casos en que los tags vienen escapados como &lt;p&gt; (ej: feeds RSS de WWR).
func CleanHTML(desc string) string {
	desc = strings.TrimSpace(desc)

	// 1. Decodificar entities nombradas PRIMERO (antes de procesar tags),
	//    para que &lt;p&gt; se convierta en <p> y pueda ser eliminado.
	desc = htmlEntities.Replace(desc)

	// 2. Eliminar HTML entities numéricas (&#160; &#x00A0; etc.)
	desc = reEntityNum.ReplaceAllString(desc, " ")

	// 3. Eliminar bloques <script> y <style> completos
	desc = reScript.ReplaceAllString(desc, "")
	desc = reStyle.ReplaceAllString(desc, "")

	// 4. Eliminar comentarios HTML <!-- ... -->
	desc = reComment.ReplaceAllString(desc, "")

	// 5. Convertir \r\n a \n para unificar saltos de línea originales
	desc = strings.ReplaceAll(desc, "\r\n", "\n")

	// 6. Convertir tags de bloque a salto de línea (preserva estructura)
	desc = reBlockTag.ReplaceAllString(desc, "\n")

	// 7. Eliminar todos los tags HTML restantes (inline: span, a, strong, etc.)
	desc = reHTMLTag.ReplaceAllString(desc, "")

	// 8. Eliminar URLs sueltas que hayan quedado en el texto plano
	desc = reURL.ReplaceAllString(desc, "")

	// 9. Normalizar espacios/tabs múltiples dentro de cada línea
	desc = strings.ReplaceAll(desc, "\t", " ")
	desc = reMultiSpace.ReplaceAllString(desc, " ")

	// 10. Colapsar más de 2 saltos de línea consecutivos
	desc = reMultiNewline.ReplaceAllString(desc, "\n\n")

	// 11. Limpiar espacios al inicio/fin de cada línea
	lines := strings.Split(desc, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}
	// Filtrar líneas vacías consecutivas ya está cubierto por reMultiNewline,
	// pero eliminamos las vacías que queden al inicio/fin del slice.
	desc = strings.Join(lines, "\n")
	desc = strings.TrimSpace(desc)

	// 12. Convertir saltos de línea reales a literal \n para no romper el CSV.
	desc = strings.ReplaceAll(desc, "\n", `\n`)

	return desc
}

func StateToISO(text string) string {
	text = strings.ToLower(strings.TrimSpace(text))
	if text == "" {
		return ""
	}

	if strings.Contains(text, ",") {
		parts := strings.Split(text, ",")
		for _, p := range parts {
			iso := StateToISO(strings.TrimSpace(p))
			if iso != "" {
				return iso
			}
		}
	}

	switch {
	case strings.Contains(text, "alabama"), strings.Contains(text, "alaska"), strings.Contains(text, "arizona"),
		strings.Contains(text, "arkansas"), strings.Contains(text, "california"), strings.Contains(text, "colorado"),
		strings.Contains(text, "connecticut"), strings.Contains(text, "delaware"), strings.Contains(text, "florida"),
		strings.Contains(text, "georgia"), strings.Contains(text, "hawaii"), strings.Contains(text, "idaho"),
		strings.Contains(text, "illinois"), strings.Contains(text, "indiana"), strings.Contains(text, "iowa"),
		strings.Contains(text, "kansas"), strings.Contains(text, "kentucky"), strings.Contains(text, "louisiana"),
		strings.Contains(text, "maine"), strings.Contains(text, "maryland"), strings.Contains(text, "massachusetts"),
		strings.Contains(text, "michigan"), strings.Contains(text, "minnesota"), strings.Contains(text, "mississippi"),
		strings.Contains(text, "missouri"), strings.Contains(text, "montana"), strings.Contains(text, "nebraska"),
		strings.Contains(text, "nevada"), strings.Contains(text, "new hampshire"), strings.Contains(text, "new jersey"),
		strings.Contains(text, "new mexico"), strings.Contains(text, "new york"), strings.Contains(text, "nueva york"),
		strings.Contains(text, "north carolina"), strings.Contains(text, "north dakota"), strings.Contains(text, "ohio"),
		strings.Contains(text, "oklahoma"), strings.Contains(text, "oregon"), strings.Contains(text, "pennsylvania"),
		strings.Contains(text, "rhode island"), strings.Contains(text, "south carolina"), strings.Contains(text, "south dakota"),
		strings.Contains(text, "tennessee"), strings.Contains(text, "texas"), strings.Contains(text, "utah"),
		strings.Contains(text, "vermont"), strings.Contains(text, "virginia"), strings.Contains(text, "washington"),
		strings.Contains(text, "west virginia"), strings.Contains(text, "wisconsin"), strings.Contains(text, "wyoming"),
		strings.Contains(text, "estados unidos"), strings.Contains(text, "ee. uu."), strings.Contains(text, "ee uu"),
		strings.Contains(text, "usa"), strings.Contains(text, "austin"), strings.Contains(text, "los ángeles"),
		strings.Contains(text, "los angeles"), strings.Contains(text, "eden prairie"), strings.Contains(text, "jersey"),
		strings.Contains(text, " dallas"), strings.Contains(text, " houston"), strings.Contains(text, " chicago"),
		strings.Contains(text, " atlanta"), strings.Contains(text, " phoenix"), strings.Contains(text, " boston"),
		strings.Contains(text, " denver"), strings.Contains(text, " seattle"), strings.Contains(text, " san francisco"),
		strings.Contains(text, ", tx"), strings.Contains(text, ", ca"), strings.Contains(text, ", ny"), strings.Contains(text, ", fl"),
		strings.Contains(text, ", ga"), strings.Contains(text, ", il"), strings.Contains(text, ", pa"), strings.Contains(text, ", oh"),
		strings.Contains(text, ", mi"), strings.Contains(text, ", nc"), strings.Contains(text, ", nj"), strings.Contains(text, ", va"),
		strings.Contains(text, ", co"), strings.Contains(text, ", wa"), strings.Contains(text, ", ma"), strings.Contains(text, ", md"),
		strings.Contains(text, ", or"), strings.Contains(text, ", az"), strings.Contains(text, ", ut"), strings.Contains(text, ", mi"),
		strings.Contains(text, ", mn"), strings.Contains(text, ", wi"), strings.Contains(text, ", tn"), strings.Contains(text, ", sc"),
		strings.Contains(text, ", ky"), strings.Contains(text, ", nv"), strings.Contains(text, ", or"), strings.Contains(text, ", ok"),
		strings.Contains(text, "remote-us"), strings.Contains(text, "remote us"), strings.Contains(text, "remote - us"):
		return "US"
	case strings.Contains(text, "ontario"), strings.Contains(text, "columbia británica"), strings.Contains(text, "british columbia"),
		strings.Contains(text, "vancouver"), strings.Contains(text, "calgary"), strings.Contains(text, "montreal"),
		strings.Contains(text, "canadá"), strings.Contains(text, "canada"):
		return "CA"
	case strings.Contains(text, "londres"), strings.Contains(text, "london"), strings.Contains(text, "reino unido"), strings.Contains(text, "uk"),
		strings.Contains(text, "manchester"), strings.Contains(text, "birmingham"), strings.Contains(text, "edinburgh"):
		return "GB"
	case strings.Contains(text, "ámsterdam"), strings.Contains(text, "amsterdam"), strings.Contains(text, "países bajos"), strings.Contains(text, "netherlands"), strings.Contains(text, "la haya"), strings.Contains(text, "the hague"), strings.Contains(text, "neerlandés"), strings.Contains(text, "dutch"):
		return "NL"
	case strings.Contains(text, "brno"), strings.Contains(text, "pilsenský kraj"), strings.Contains(text, "república checa"), strings.Contains(text, "czech"):
		return "CZ"
	case strings.Contains(text, "zug"), strings.Contains(text, "suiza"), strings.Contains(text, "switzerland"):
		return "CH"
	case strings.Contains(text, "irlanda"), strings.Contains(text, "ireland"), strings.Contains(text, "dublín"), strings.Contains(text, "dublin"):
		return "IE"
	case strings.Contains(text, "malta"), strings.Contains(text, "santa venera"):
		return "MT"
	case strings.Contains(text, "austria"), strings.Contains(text, "klagenfurt"), strings.Contains(text, "graz"), strings.Contains(text, "viena"), strings.Contains(text, "vienna"):
		return "AT"
	case strings.Contains(text, "ucrania"), strings.Contains(text, "ukraine"), strings.Contains(text, "kiev"), strings.Contains(text, "kyiv"):
		return "UA"
	case strings.Contains(text, "andorra"):
		return "AD"
	case strings.Contains(text, "alemania"), strings.Contains(text, "germany"), strings.Contains(text, "berlín"), strings.Contains(text, "berlin"), strings.Contains(text, "múnich"), strings.Contains(text, "munich"),
		strings.Contains(text, "gmbh"), strings.Contains(text, "analogue insight"):
		return "DE"
	case strings.Contains(text, "francia"), strings.Contains(text, "france"), strings.Contains(text, "parís"), strings.Contains(text, "paris"),
		strings.Contains(text, "eviden"):
		return "FR"
	case strings.Contains(text, "españa"), strings.Contains(text, "spain"), strings.Contains(text, "madrid"), strings.Contains(text, "barcelona"):
		return "ES"
	case strings.Contains(text, "italia"), strings.Contains(text, "italy"), strings.Contains(text, "roma"), strings.Contains(text, "rome"):
		return "IT"
	case strings.Contains(text, "mexico"), strings.Contains(text, "méxico"):
		return "MX"
	case strings.Contains(text, "argentina"), strings.Contains(text, "buenos aires"):
		return "AR"
	case strings.Contains(text, "brasil"), strings.Contains(text, "brazil"), strings.Contains(text, "são paulo"),
		strings.Contains(text, " br "), strings.Contains(text, ": br"), strings.Contains(text, ", br"):
		return "BR"
	case strings.Contains(text, "filipinas"), strings.Contains(text, "philippines"), strings.Contains(text, "manila"),
		strings.Contains(text, "cagayan de oro"):
		return "PH"
	case strings.Contains(text, "israel"), strings.Contains(text, "tel aviv"), strings.Contains(text, "jerusalén"), strings.Contains(text, "jerusalem"):
		return "IL"
	case strings.Contains(text, "india"), strings.Contains(text, "bangalore"), strings.Contains(text, "bengaluru"), strings.Contains(text, "mumbai"), strings.Contains(text, "delhi"),
		strings.Contains(text, "odisha"), strings.Contains(text, "uttar pradesh"):
		return "IN"
	case strings.Contains(text, "honduras"), strings.Contains(text, "tegucigalpa"), strings.Contains(text, "francisco morazán"):
		return "HN"
	case strings.Contains(text, "guatemala"):
		return "GT"
	case strings.Contains(text, "nicaragua"):
		return "NI"
	case strings.Contains(text, "costa rica"):
		return "CR"
	case strings.Contains(text, "panamá"), strings.Contains(text, "panama"):
		return "PA"
	case strings.Contains(text, "el salvador"):
		return "SV"
	case strings.Contains(text, "república dominicana"), strings.Contains(text, "dominican republic"):
		return "DO"
	case strings.Contains(text, "ecuador"), strings.Contains(text, "quito"), strings.Contains(text, "guayaquil"):
		return "EC"
	case strings.Contains(text, "venezuela"), strings.Contains(text, "caracas"):
		return "VE"
	case strings.Contains(text, "perú"), strings.Contains(text, "peru"), strings.Contains(text, "lima"):
		return "PE"
	case strings.Contains(text, "bolivia"), strings.Contains(text, "la paz"):
		return "BO"
	case strings.Contains(text, "paraguay"), strings.Contains(text, "asunción"):
		return "PY"
	case strings.Contains(text, "puerto rico"):
		return "PR"
	case strings.Contains(text, "singapur"), strings.Contains(text, "singapore"):
		return "SG"
	case strings.Contains(text, "japón"), strings.Contains(text, "japan"), strings.Contains(text, "tokio"), strings.Contains(text, "tokyo"):
		return "JP"
	case strings.Contains(text, "australia"), strings.Contains(text, "sydney"), strings.Contains(text, "melbourne"):
		return "AU"
	case strings.Contains(text, "nueva zelanda"), strings.Contains(text, "new zealand"), strings.Contains(text, "auckland"):
		return "NZ"
	case strings.Contains(text, "chile"), strings.Contains(text, "santiago"):
		return "CL"
	case strings.Contains(text, "colombia"), strings.Contains(text, "bogotá"), strings.Contains(text, "medellín"):
		return "CO"
	case strings.Contains(text, "uruguay"), strings.Contains(text, "montevideo"):
		return "UY"
	case strings.Contains(text, "estonia"), strings.Contains(text, "tallinn"), strings.Contains(text, "tallin"):
		return "EE"
	case strings.Contains(text, "polonia"), strings.Contains(text, "poland"), strings.Contains(text, "warsaw"), strings.Contains(text, "varsovia"):
		return "PL"
	case strings.Contains(text, "bélgica"), strings.Contains(text, "belgium"), strings.Contains(text, "brussels"), strings.Contains(text, "bruselas"):
		return "BE"
	case strings.Contains(text, "finlandia"), strings.Contains(text, "finland"), strings.Contains(text, "helsinki"):
		return "FI"
	case strings.Contains(text, "sudáfrica"), strings.Contains(text, "south africa"), strings.Contains(text, "cape town"), strings.Contains(text, "johannesburg"):
		return "ZA"
	case strings.Contains(text, "líbano"), strings.Contains(text, "lebanon"):
		return "LB"
	case strings.Contains(text, "dinamarca"), strings.Contains(text, "denmark"), strings.Contains(text, "creativeforce.team"):
		return "DK"
	case strings.Contains(text, "chipre"), strings.Contains(text, "cyprus"), strings.Contains(text, "paphos"):
		return "CY"
	case strings.Contains(text, "european union"), strings.Contains(text, "unión europea"), strings.Contains(text, " eu "):
		return "EU"
	case strings.Contains(text, "tailandia"), strings.Contains(text, "thailand"):
		return "TH"
	case strings.Contains(text, "indonesia"), strings.Contains(text, "bali"), strings.Contains(text, "jakarta"):
		return "ID"
	case strings.Contains(text, "vietnam"), strings.Contains(text, "viet nam"):
		return "VN"
	case strings.Contains(text, "egipto"), strings.Contains(text, "egypt"):
		return "EG"
	case strings.Contains(text, "china"), strings.Contains(text, "bamboo-works.com"):
		return "CN"
	case strings.Contains(text, "north america"), strings.Contains(text, "norteamérica"):
		return "US"
	case strings.Contains(text, "wolfsburg"):
		return "DE"
	case strings.Contains(text, "redwood city"), strings.Contains(text, "san francisco"), strings.Contains(text, "seattle"), strings.Contains(text, "salt lake city"), strings.Contains(text, "dallas"), strings.Contains(text, "new york"):
		return "US"
	case strings.Contains(text, "scoutsolutions.net"), strings.Contains(text, "quinstreet.com"), strings.Contains(text, "theclassconsultinggroup.org"), strings.Contains(text, "intro.io"), strings.Contains(text, "quinstreet"),
		strings.Contains(text, "amgen"):
		return "US"
	case strings.Contains(text, "codersbrain.com"):
		return "IN"
	case strings.Contains(text, "canonical"), strings.Contains(text, "adria solutions"), strings.Contains(text, "trt solutions"):
		return "GB"
	case strings.Contains(text, "serbia"), strings.Contains(text, "belgrade"), strings.Contains(text, "belgrado"):
		return "RS"
	case strings.Contains(text, "grecia"), strings.Contains(text, "greece"), strings.Contains(text, "atenas"), strings.Contains(text, "athens"):
		return "GR"
	case strings.Contains(text, "turquía"), strings.Contains(text, "turkey"), strings.Contains(text, "estambul"), strings.Contains(text, "istanbul"):
		return "TR"
	case strings.Contains(text, "united states"), strings.Contains(text, "u.s.a"), strings.Contains(text, "u.s."), strings.Contains(text, "ee.uu"), strings.Contains(text, "eeuu"), strings.Contains(text, "estados unidos"):
		return "US"
	case strings.Contains(text, "united kingdom"), strings.Contains(text, "reino unido"), strings.Contains(text, "uk"):
		return "GB"
	case strings.Contains(text, "suecia"), strings.Contains(text, "sweden"):
		return "SE"
	case strings.Contains(text, "noruega"), strings.Contains(text, "norway"):
		return "NO"
	}
	return ""
}

func ExtractISO(content string) string {
	for i := 0; i < len(content)-3; {
		r1 := rune(content[i])
		if r1 >= 0x80 {
			r, size := utf8.DecodeRuneInString(content[i:])
			if r >= 0x1F1E6 && r <= 0x1F1FF {
				first := string(r - 0x1F1E6 + 'A')
				i += size
				if i < len(content) {
					r2, size2 := utf8.DecodeRuneInString(content[i:])
					if r2 >= 0x1F1E6 && r2 <= 0x1F1FF {
						second := string(r2 - 0x1F1E6 + 'A')
						return first + second
					}
					i += size2
					continue
				}
			}
			i += size
		} else {
			i++
		}
	}
	return ""
}

func ParseDate(dateStr string) string {
	if dateStr == "" {
		return ""
	}

	replacer := strings.NewReplacer(
		"Lun,", "Mon,", "Mar,", "Tue,", "Mié,", "Wed,", "Jue,", "Thu,", "Vie,", "Fri,", "Sáb,", "Sat,", "Dom,", "Sun,",
		"Ene", "Jan", "Abr", "Apr", "Ago", "Aug", "Dic", "Dec",
	)

	englishDate := replacer.Replace(dateStr)

	layouts := []string{
		"Mon, 02 Jan 2006 15:04:05 -0700",
		"Mon, 02 Jan 2006 15:04:05 GMT",
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02 15:04:05",
	}

	for _, layout := range layouts {
		t, err := time.Parse(layout, englishDate)
		if err == nil {
			return t.Format("2006-01-02 15:04:05")
		}
	}

	return dateStr
}

func FetchURL(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 429 {
		return "", fmt.Errorf("RATE_LIMIT")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
