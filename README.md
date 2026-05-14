# JobAR (Jobs + API + RSS)

This application aggregates job listings from various platforms. It fetches, parses, and processes job postings from multiple providers to standardize the data.

## Types of Providers

The application currently supports two main types of job data providers: **RSS Feeds** and **JSON APIs**.

### 1. RSS Feed Providers
These providers expose their job listings via standard RSS/XML feeds. The application parses the XML `<item>` blocks to extract details like title, description, company, location, and link.

- **Remotive**
- **WeWorkRemotely**
- **FindJobIt**
- **RemoteOK**
- **Jobicy**

### 2. JSON API Providers
These providers expose a RESTful API returning data in JSON format. The application fetches the JSON and unmarshals it into internal structures.

- **Himalayas**
