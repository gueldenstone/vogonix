export namespace jira {
	
	export class Issue {
	    Summary: string;
	    Assignee: string;
	    Key: string;
	
	    static createFrom(source: any = {}) {
	        return new Issue(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Summary = source["Summary"];
	        this.Assignee = source["Assignee"];
	        this.Key = source["Key"];
	    }
	}
	export class Worklog {
	    Duration: number;
	    Comment: string;
	
	    static createFrom(source: any = {}) {
	        return new Worklog(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Duration = source["Duration"];
	        this.Comment = source["Comment"];
	    }
	}

}

