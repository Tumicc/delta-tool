export namespace app {
	
	export class WeaponCode {
	    id: string;
	    mode: string;
	    name: string;
	    tier: string;
	    price?: number;
	    build: string;
	    code: string;
	    range?: number;
	    update_time?: string;
	    source: string;
	
	    static createFrom(source: any = {}) {
	        return new WeaponCode(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.mode = source["mode"];
	        this.name = source["name"];
	        this.tier = source["tier"];
	        this.price = source["price"];
	        this.build = source["build"];
	        this.code = source["code"];
	        this.range = source["range"];
	        this.update_time = source["update_time"];
	        this.source = source["source"];
	    }
	}

}

