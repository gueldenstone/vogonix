// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {jira} from '../models';

export function GetAssignedIssues():Promise<Array<jira.Issue>>;

export function GetBaseUrl():Promise<string>;

export function GetCurrentTimerValue(arg1:string):Promise<number>;

export function GetTimeSpentOnIssue(arg1:string):Promise<string>;

export function GetWorkLogs(arg1:string):Promise<Array<jira.Worklog>>;

export function LogDebugf(arg1:string,arg2:Array<any>):Promise<void>;

export function LogWarning(arg1:string):Promise<void>;

export function LogWarningf(arg1:string,arg2:Array<any>):Promise<void>;

export function PauseTimer(arg1:string):Promise<void>;

export function ResetTimer(arg1:string):Promise<void>;

export function StartTimer(arg1:string):Promise<void>;

export function SubmitWorklog(arg1:string,arg2:number):Promise<void>;
