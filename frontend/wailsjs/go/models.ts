export namespace main {
	
	export class DeviceInfo {
	    model: string;
	    firmware: string;
	    battery: number;
	    anc: string;
	
	    static createFrom(source: any = {}) {
	        return new DeviceInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.model = source["model"];
	        this.firmware = source["firmware"];
	        this.battery = source["battery"];
	        this.anc = source["anc"];
	    }
	}

}

