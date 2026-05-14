package tracks

type Track struct {
	Name        string
	Description string
}

type Area struct {
	ID     string
	Name   string
	Tracks []Track
}

var AllAreas = []Area{
	{
		ID:   "AGILE",
		Name: "AGILE",
		Tracks: []Track{
			{Name: "Scrum Master", Description: "Para impulsar la adopción de metodologías ágiles."},
			{Name: "Agile Coach", Description: "Para acompañar a equipos y líderes en transformación ágil."},
		},
	},
	{
		ID:   "ANALIST",
		Name: "Analisis de Negocios",
		Tracks: []Track{
			{Name: "Functional Analyst", Description: "Para traducir necesidades de negocio en requerimientos."},
			{Name: "Business Analyst", Description: "Para analizar procesos e implementar mejoras en productos."},
			{Name: "Product Business Analyst", Description: "Conecta negocio, datos y tecnología para impulsar productos. "},
		},
	},
	{
		ID:   "ANDROID",
		Name: "Android",
		Tracks: []Track{
			{Name: "Android Programmer", Description: "La base para desarrollar aplicaciones móviles en Android."},
			{Name: "Associate Android Developer", Description: "Para desarrollar apps Android listas para producción real."},
		},
	},
	{
		ID:   "CLOUD",
		Name: "Cloud",
		Tracks: []Track{
			{Name: "Cloud Support Technician", Description: "Administra servicios e infraestructura en la nube de AWS."},
			{Name: "Cloud Engineer", Description: "Para diseñar e implementar soluciones escalables cloud."},
			{Name: "Cloud Solution Architect", Description: "Diseña soluciones seguras para infraestructura cloud."},
		},
	},
	{
		ID:   "DATA",
		Name: "DATA SCIENCE",
		Tracks: []Track{
			{Name: "Data Analyst", Description: "Para analizar datos y apoyar la toma de decisiones."},
			{Name: "Data Engineer", Description: "Para construir pipelines de datos para análisis y BI."},
			{Name: "Data Scientist", Description: "Para aplicar estadística y programación sobre datos."},
			{Name: "MLOps Engineer", Description: "Para implementar y gestionar modelos de IA en producción."},
		},
	},
	{
		ID:   "DEVOPS",
		Name: "Devops",
		Tracks: []Track{
			{Name: "Devops", Description: "Automatiza procesos entre desarrollo y operaciones."},
			{Name: "DevSecOps", Description: "Implementa pipelines de DevOps y prácticas de seguridad."},
			{Name: "Site Reliability Engineer (SRE)", Description: "Garantiza estabilidad de sistemas a gran escala."},
		},
	},
	{
		ID:   "DWEB",
		Name: "Diseño UI/WEB",
		Tracks: []Track{
			{Name: "UI Designer", Description: "Para diseñar interfaces centradas en la experiencia del usuario."},
			{Name: "Web Designer", Description: "Para transformar ideas en sitios web efectivos y profesionales."},
			{Name: "Web Designer Certified", Description: "Para validar conocimientos con certificación internacional."},
		},
	},
	{
		ID:   "ECOMM",
		Name: "eCommerce",
		Tracks: []Track{
			{Name: "SEO Analyst", Description: "Optimiza sitios para posicionamiento en buscadores y LLM’s."},
			{Name: "E-Commerce Manager", Description: "Para administrar plataformas de eCommerce y canales."},
			{Name: "eCommerce Data Analyst", Description: "Para administrar eCommerce basado en analítica de datos."},
		},
	},
	{
		ID:   "FEND",
		Name: "Frontend Developer",
		Tracks: []Track{
			{Name: "Front End Developer", Description: "Todo lo necesario para crear interfaces web funcionales."},
			{Name: "React JS Developer", Description: "Para trabajar en el framework más demandado para apps dinámicas."},
			{Name: "Front End Developer Certified", Description: "Para validar conocimientos con certificación internacional."},
		},
	},
	{
		ID:   "IA",
		Name: "IA",
		Tracks: []Track{
			{Name: "Automation Specialist", Description: "Para crear agentes y automatizar tareas y procesos."},
			{Name: "AI Consultant", Description: "Para liderar procesos implementando herramientas de IA."},
			{Name: "AI Strategy Consultant", Description: "Para validar conocimientos con certificación internacional."},
		},
	},
	{
		ID:   "JAVA",
		Name: "JAVA",
		Tracks: []Track{
			{Name: "Backend Java Developer", Description: "Todo lo necesario para iniciar como desarrollador backend con Java."},
			{Name: "Java Software Engineer", Description: "Para diseñar y construir soluciones avanzadas con Java."},
			{Name: "Java Enterprise Certified", Description: "Para validar conocimientos con certificación internacional."},
		},
	},
	{
		ID:   "JS",
		Name: "JAVASCRIPT",
		Tracks: []Track{
			{Name: "Javascript Developer", Description: "La base ideal para iniciar tu carrera en desarrollo con JS."},
			{Name: "Javascript Fullstack Developer", Description: "Para convertirte en desarrollador front-end y back-end."},
			{Name: "Certified Associate JavaScript Programmer", Description: "Para validar conocimientos con certificación internacional."},
		},
	},
	{
		ID:   "LNX",
		Name: "Linux",
		Tracks: []Track{
			{Name: "Linux Administrator", Description: "Para administrar sistemas operativos GNU/Linux."},
			{Name: "Linux Engineer", Description: "Diseña e implementa soluciones basadas en Linux."},
			{Name: "Linux Platform Engineer", Description: "Diseña plataformas Linux para entornos empresariales."},
		},
	},
	{
		ID:   "MKO",
		Name: "Marketing Digital",
		Tracks: []Track{
			{Name: "Community Manager", Description: "Para gestionar redes sociales y comunidades de marcas."},
			{Name: "Paid Media Manager", Description: "Para planificar y gestionar publicidad digital."},
			{Name: "Digital Marketing Analyst", Description: "Diseña campañas orientadas a crecimiento y marca."},
			{Name: "Growth Hacking Analyst", Description: "Impulsa el crecimiento mediante tácticas innovadoras."},
		},
	},
	{
		ID:   "MULT",
		Name: "Multimedia",
		Tracks: []Track{
			{Name: "CAD Designer", Description: "Crea planos y modelos digitales para diseño técnico."},
			{Name: "Motion Designer", Description: "Diseña contenido audiovisual con efectos y motion graphics."},
			{Name: "Multimedia Designer", Description: "Diseña contenido multimedia con herramientas creativas."},
		},
	},
	{
		ID:   "NET",
		Name: ".NET",
		Tracks: []Track{
			{Name: "Backend NET Developer", Description: "Para desarrollar aplicaciones backend con Microsoft.NET."},
			{Name: "NET Software Engineer", Description: "Para llevar tus habilidades a nivel ingeniería y software escalable."},
			{Name: "Azure Solutions Developer", Description: "Para validar conocimientos con certificación internacional."},
		},
	},
	{
		ID:   "OFFICE",
		Name: "OFFICE",
		Tracks: []Track{
			{Name: "Administrative Assistant", Description: "Gestiona tareas administrativas con Office."},
			{Name: "Microsoft Office Specialist", Description: "Administra planillas complejas y reportes financieros."},
			{Name: "Power BI Developer", Description: "Para validar conocimientos con certificación internacional."},
		},
	},
	{
		ID:   "ORACLE",
		Name: "ORACLE",
		Tracks: []Track{
			{Name: "PL/SQL Developer", Description: "Crea consultas y lógica de negocio en Oracle."},
			{Name: "Oracle DBA", Description: "Administra bases de datos Oracle empresariales."},
			{Name: "Oracle Database Engineer", Description: "Diseña y optimiza la arquitectura de bases de datos."},
		},
	},
	{
		ID:   "PHP",
		Name: "PHP",
		Tracks: []Track{
			{Name: "PHP Developer", Description: "Para comenzar a desarrollar aplicaciones web con PHP."},
			{Name: "PHP Fullstack Developer", Description: "Para convertirte en desarrollador front-end y back-end."},
			{Name: "PHP Certified Engineer", Description: "Para validar conocimientos con certificación internacional."},
		},
	},
	{
		ID:   "PRODUCT",
		Name: "Producto",
		Tracks: []Track{
			{Name: "Product Owner (PO)", Description: "Conecta necesidades del negocio con desarrollo de producto."},
			{Name: "Product Manager", Description: "Lidera equipos y prioridades para el desarrollo del producto."},
			{Name: "AI Product Manager", Description: "Define estrategias de producto basadas en IA y analítica."},
		},
	},
	{
		ID:   "PROJECT",
		Name: "PROJECT",
		Tracks: []Track{
			{Name: "Project Manager", Description: "Gestiona equipos y recursos para proyectos."},
			{Name: "IT Project Manager", Description: "Lidera proyectos tecnológicos en organizaciones."},
			{Name: "Agile Project Manager", Description: "Para validar conocimientos con certificación internacional."},
		},
	},
	{
		ID:   "PYTHON",
		Name: "Python",
		Tracks: []Track{
			{Name: "Python Developer", Description: "Todo lo que necesitas para empezar como desarrollador Python."},
			{Name: "Full Stack Python", Description: "Para convertirte en un desarrollador de Front-End y Back-End."},
			{Name: "Python Engineer", Description: "Para convertirte en un experto en Python."},
			{Name: "Certified Associate Python Programmer", Description: "Para validar conocimientos con certificación internacional."},
		},
	},
	{
		ID:   "REDES",
		Name: "Redes",
		Tracks: []Track{
			{Name: "Network Analyst", Description: "Implementa y analiza el funcionamiento de redes e infra."},
			{Name: "Network Administrator", Description: "Administra redes corporativas e infraestructuras."},
			{Name: "Network Engineer", Description: "Desarrolla infraestructuras de red empresariales."},
		},
	},
	{
		ID:   "ROB",
		Name: "Robótica",
		Tracks: []Track{
			{Name: "AI Application Developer", Description: "Para desarrollar apps que integran IA en soluciones reales."},
			{Name: "Robotics AI Engineer", Description: "Para desarrollar soluciones de software y hardware con robótica."},
			{Name: "Artificial Intelligence Professional Certificate", Description: "Para validar conocimientos con certificación internacional."},
		},
	},
	{
		ID:   "SINF",
		Name: "Seguridad Inf.",
		Tracks: []Track{
			{Name: "SOC Analyst", Description: "Monitorea y detecta incidentes en tiempo real."},
			{Name: "Penetration Tester", Description: "Ethical hacking para detectar vulnerabilidades."},
			{Name: "Security Engineer", Description: "Implementa blue y red teams para seguridad IT."},
			{Name: "Enterprise Security Architect", Description: "Lidera la visión estratégica de ciberseguridad."},
		},
	},
	{
		ID:   "SOPORTE",
		Name: "Soporte Técnico",
		Tracks: []Track{
			{Name: "Help Desk Agent", Description: "Brinda asistencia técnica y resuelve incidencias básicas."},
			{Name: "IT Support Specialist", Description: "Gestiona incidencias técnicas y soporte nivel 1 y 2."},
			{Name: "IT Support Manager", Description: "Coordina la operación y mejora continua del soporte IT."},
		},
	},
	{
		ID:   "SQL",
		Name: "SQL Server",
		Tracks: []Track{
			{Name: "SQL Developer", Description: "Diseña consultas, reportes y lógica de datos. "},
			{Name: "SQL Database Administrator (DBA)", Description: "Administra bases de datos SQL empresariales."},
			{Name: "Azure Database Administrator Associate", Description: "Para validar conocimientos con certificación internacional."},
		},
	},
	{
		ID:   "TALENTO",
		Name: "Talento",
		Tracks: []Track{
			{Name: "IT Recruiter", Description: "Para atraer profesionales tecnológicos en organizaciones."},
			{Name: "HR Data Analyst", Description: "Impulsa decisiones estratégicas de RR. HH. con analítica."},
			{Name: "People Manager", Description: "Gestiona equipos y desarrolla el talento organizacional."},
			{Name: "Agile HR Certified Professional", Description: "Para validar conocimientos con certificación internacional."},
		},
	},
	{
		ID:   "TQA",
		Name: "TQA",
		Tracks: []Track{
			{Name: "QA Analyst", Description: "La base ideal para comenzar como tester manual."},
			{Name: "QA Automation", Description: "Para especializarte en pruebas automatizadas y optimización."},
			{Name: "QA Manager", Description: "Para liderar equipos y estrategias de calidad de software."},
			{Name: "Tester Foundation Level Certified", Description: "Para validar conocimientos con certificación internacional."},
		},
	},
	{
		ID:   "UX",
		Name: "UX",
		Tracks: []Track{
			{Name: "UX Researcher", Description: "Para investigar comportamiento y mejorar experiencia digital."},
			{Name: "UX/UI Designer", Description: "Para diseñar productos que combinan usabilidad y estética."},
			{Name: "UX Manager", Description: "Para liderar equipos de contenido, diseño e investigación. "},
		},
	},
	{
		ID:   "WS",
		Name: "Windows Server",
		Tracks: []Track{
			{Name: "Windows System Administrator", Description: "Administra servidores y servicios en entornos Windows."},
			{Name: "Microsoft Systems Engineer", Description: "Administra infraestructura híbrida local y cloud."},
			{Name: "Microsoft Infrastructure Architect", Description: "Diseña soluciones escalables para entornos Microsoft."},
		},
	},
}
