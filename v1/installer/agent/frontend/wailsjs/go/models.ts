export namespace main {
	
	export class AppConfig {
	    AppName: string;
	    Description: string;
	    Icon: string;
	    Tagline: string;
	    Title: string;
	    Width: number;
	    Height: number;
	    Min: boolean;
	    Max: boolean;
	    Quit: boolean;
	    Version: string;
	    Domain: string;
	    installDir: string;
	
	    static createFrom(source: any = {}) {
	        return new AppConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.AppName = source["AppName"];
	        this.Description = source["Description"];
	        this.Icon = source["Icon"];
	        this.Tagline = source["Tagline"];
	        this.Title = source["Title"];
	        this.Width = source["Width"];
	        this.Height = source["Height"];
	        this.Min = source["Min"];
	        this.Max = source["Max"];
	        this.Quit = source["Quit"];
	        this.Version = source["Version"];
	        this.Domain = source["Domain"];
	        this.installDir = source["installDir"];
	    }
	}

}

