export namespace main {
	
	export class Capabilities {
	    hasANC: boolean;
	    hasEQ: boolean;
	    hasBattery: boolean;
	    hasVolume: boolean;
	    hasRoomPlacement: boolean;
	    hasPartyMode: boolean;
	    hasNightMode: boolean;
	    hasLED: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Capabilities(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.hasANC = source["hasANC"];
	        this.hasEQ = source["hasEQ"];
	        this.hasBattery = source["hasBattery"];
	        this.hasVolume = source["hasVolume"];
	        this.hasRoomPlacement = source["hasRoomPlacement"];
	        this.hasPartyMode = source["hasPartyMode"];
	        this.hasNightMode = source["hasNightMode"];
	        this.hasLED = source["hasLED"];
	    }
	}
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
	export class ScannedDevice {
	    name: string;
	    address: string;
	
	    static createFrom(source: any = {}) {
	        return new ScannedDevice(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.address = source["address"];
	    }
	}

}

