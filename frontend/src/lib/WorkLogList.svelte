<script lang="ts">
    import {
        Accordion,
        AccordionItem,
        Listgroup,
        Spinner,
    } from "flowbite-svelte";
    import { formatDurationFromNanoseconds } from "./util";
    import * as jira from "@go/jira/JiraInstance";
    import * as wails from "@runtime/runtime";
    import * as models from "@go/models";
    import { timeAgo } from "@lib/util";
    import * as Icons from "flowbite-svelte-icons";

    // export let workLogs: models.jira.Worklog[];
    export let issueKey: string;

    function sumWorkLogs(worklogs: models.jira.Worklog[]): number {
        return worklogs.reduce((total, worklog) => total + worklog.duration, 0);
    }
    $: workLogsPromise = jira.GetWorkLogs(issueKey);

    export function refresh() {
        workLogsPromise = jira.GetWorkLogs(issueKey);
    }
    wails.EventsOn("update_worklogs_" + issueKey, refresh);

    let open = false;
</script>

<Accordion flush>
    {#await workLogsPromise}
        <AccordionItem bind:open>
            <Spinner size="10" />
        </AccordionItem>
    {:then workLogs}
        <AccordionItem bind:open paddingFlush="p-2">
            <span slot="header"
                >{formatDurationFromNanoseconds(sumWorkLogs(workLogs))}</span
            >
            <Listgroup items={workLogs} let:item>
                {#if item != undefined}
                    <div
                        class="flex flex-row h-full grow justify-between gap-10"
                    >
                        <div class="flex items-center gap-2">
                            <Icons.UserCircleSolid />
                            <span class="text-left"> {item.author} </span>
                        </div>
                        <div class="flex items-center gap-2">
                            <Icons.ClockSolid />
                            <span class="text-left">
                                {formatDurationFromNanoseconds(item.duration)}
                            </span>
                        </div>
                        <span class="ml-auto text-right">
                            {timeAgo(item.updated)}
                        </span>
                    </div>
                {/if}
            </Listgroup>
        </AccordionItem>
    {/await}
</Accordion>
