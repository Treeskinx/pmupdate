export namespace main {
	
	export class FileInput {
	    name: string;
	    size: number;
	    type: string;
	    data: string;
	
	    static createFrom(source: any = {}) {
	        return new FileInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.size = source["size"];
	        this.type = source["type"];
	        this.data = source["data"];
	    }
	}

}

