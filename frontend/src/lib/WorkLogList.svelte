<script lang="ts">
    import { Accordion, AccordionItem, Listgroup } from "flowbite-svelte";
    import * as models from "@go/models";
    import { formatDurationFromNanoseconds } from "./util";

    export let workLogs: models.jira.Worklog[];

    function sumWorkLogs(worklogs: models.jira.Worklog[]): number {
        return worklogs.reduce((total, worklog) => total + worklog.duration, 0);
    }
</script>

<Accordion flush>
    <AccordionItem>
        <span slot="header"
            >{formatDurationFromNanoseconds(sumWorkLogs(workLogs))}</span
        >
        <Listgroup items={workLogs} let:item>
            {#if item != undefined}
                {item.author} on {item.updated}: {formatDurationFromNanoseconds(
                    item.duration,
                )}
            {/if}
        </Listgroup>
    </AccordionItem>
</Accordion>
