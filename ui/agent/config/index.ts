export const URL = {
    "Domain":"fehmicorp.in",
    "Website": "https://fehmicorp.in"
}

export const Config = {
    "Name": "Fehmi",
    "Assets":{
        "logo": "/logo.png",
    },
    "Company": {
        "LegalName": "Fehmi Corporation",
    },
    "Apps": {
        "Title": "Fehmi Agent Installer",
        "Tagline":"Managing Cloud makes easy",
        "Description": "Fehmi Agent Installer is a web-based application that allows users to easily install and manage Fehmi agents on their systems. It provides a user-friendly interface for configuring and deploying agents, making it simple for users to monitor and control their cloud infrastructure.",
    
    }
} as const;

export type UrlDetails = typeof URL;

export type AppDetails = typeof Config;