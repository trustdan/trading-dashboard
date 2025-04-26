export namespace models {
	
	export class RiskAssessment {
	    id: string;
	    // Go type: time
	    date: any;
	    emotional: number;
	    fomo: number;
	    bias: number;
	    physical: number;
	    pnl: number;
	    overallScore: number;
	
	    static createFrom(source: any = {}) {
	        return new RiskAssessment(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.date = this.convertValues(source["date"], null);
	        this.emotional = source["emotional"];
	        this.fomo = source["fomo"];
	        this.bias = source["bias"];
	        this.physical = source["physical"];
	        this.pnl = source["pnl"];
	        this.overallScore = source["overallScore"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class StockRating {
	    id: string;
	    // Go type: time
	    date: any;
	    ticker: string;
	    marketSentiment: number;
	    basicMaterials: number;
	    communicationServices: number;
	    consumerCyclical: number;
	    consumerDefensive: number;
	    energy: number;
	    financial: number;
	    healthcare: number;
	    industrials: number;
	    realEstate: number;
	    technology: number;
	    utilities: number;
	    stockSentiment: number;
	    pattern: string;
	    enthusiasmRating: number;
	
	    static createFrom(source: any = {}) {
	        return new StockRating(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.date = this.convertValues(source["date"], null);
	        this.ticker = source["ticker"];
	        this.marketSentiment = source["marketSentiment"];
	        this.basicMaterials = source["basicMaterials"];
	        this.communicationServices = source["communicationServices"];
	        this.consumerCyclical = source["consumerCyclical"];
	        this.consumerDefensive = source["consumerDefensive"];
	        this.energy = source["energy"];
	        this.financial = source["financial"];
	        this.healthcare = source["healthcare"];
	        this.industrials = source["industrials"];
	        this.realEstate = source["realEstate"];
	        this.technology = source["technology"];
	        this.utilities = source["utilities"];
	        this.stockSentiment = source["stockSentiment"];
	        this.pattern = source["pattern"];
	        this.enthusiasmRating = source["enthusiasmRating"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Trade {
	    id: string;
	    // Go type: time
	    entryDate: any;
	    ticker: string;
	    sector: string;
	    entryPrice: number;
	    notes: string;
	    // Go type: time
	    expirationDate: any;
	    strategyType: string;
	    spreadType: string;
	    direction: string;
	
	    static createFrom(source: any = {}) {
	        return new Trade(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.entryDate = this.convertValues(source["entryDate"], null);
	        this.ticker = source["ticker"];
	        this.sector = source["sector"];
	        this.entryPrice = source["entryPrice"];
	        this.notes = source["notes"];
	        this.expirationDate = this.convertValues(source["expirationDate"], null);
	        this.strategyType = source["strategyType"];
	        this.spreadType = source["spreadType"];
	        this.direction = source["direction"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

