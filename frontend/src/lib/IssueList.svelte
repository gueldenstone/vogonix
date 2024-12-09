<script lang="ts">
    import { Button, Card, Spinner } from "flowbite-svelte";
    import * as jira from "@go/jira/JiraInstance";
    import * as wails from "@runtime/runtime";

    import { onMount } from "svelte";
    import Issue from "./Issue.svelte";

    // local state for timer values
    $: timerValues = new Map<string, number>();

    // reactive issue promise for rendering the issues
    $: issuesPromise = jira.GetAssignedIssues();

    export function refresh() {
        issuesPromise = jira.GetAssignedIssues();
    }

    onMount(async () => {
        let issueList = await issuesPromise;
        issueList.forEach(async (issue) => {
            timerValues[issue.key] = await jira.GetCurrentTimerValue(issue.key);
        });
    });
</script>

<div class="flex p-2 gap-2 flex-col">
    {#await issuesPromise}
        <div class="text-center p-10">
            <Spinner size="10" />
        </div>
    {:then issues}
        {#each issues as issue, index (issue.key)}
            <Issue bind:timerValue={timerValues[issue.key]} {issue} />
        {/each}
    {/await}
</div>
