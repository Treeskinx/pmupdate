export namespace main {
	
	export class CSVData {
	    Dropped: string;
	
	    static createFrom(source: any = {}) {
	        return new CSVData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Dropped = source["Dropped"];
	    }
	}
	export class CSVDatas {
	    Dropped: string[];
	
	    static createFrom(source: any = {}) {
	        return new CSVDatas(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Dropped = source["Dropped"];
	    }
	}
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

