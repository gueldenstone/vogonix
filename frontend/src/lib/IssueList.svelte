<script lang="ts">
    import { Button, Card, Spinner } from "flowbite-svelte";
    import * as jira from "@go/jira/JiraInstance";
    import * as wails from "@runtime/runtime";

    import JiraIssueKey from "@lib/IssueKey.svelte";
    import StopWatch from "@lib/StopWatch.svelte";
    import * as Icons from "flowbite-svelte-icons";
    import { onMount } from "svelte";
    import WorkLogList from "./WorkLogList.svelte";

    // base url
    let jiraBaseUrl: string;
    jira.GetBaseUrl().then((url) => (jiraBaseUrl = url));

    // local state for timer values
    $: timerValues = new Map<string, number>();

    // reactive issue promise for rendering the issues
    $: issuesPromise = jira.GetAssignedIssues();

    function setupTimerEventListener(issueKey: string) {
        wails.EventsOn("timer_tick_" + issueKey, (currentTime) => {
            timerValues[issueKey] = currentTime;
        });
    }

    export function refresh() {
        issuesPromise = jira.GetAssignedIssues();
    }

    onMount(async () => {
        let issueList = await issuesPromise;
        issueList.forEach(async (issue) => {
            timerValues[issue.key] = await jira.GetCurrentTimerValue(issue.key);
        });
    });

    wails.EventsOn("update_worklogs", refresh);

    function submitPossible(issueKey: string) {
        return timerValues[issueKey] > 60 * 1e9;
    }
</script>

<div class="flex p-2 gap-2 flex-col">
    {#await issuesPromise}
        <div class="text-center p-10">
            <Spinner size="10" />
        </div>
    {:then issues}
        {#each issues as issue, index (issue.key)}
            {@const issueKey = issue.key}
            <Card size="xl" class="h-full grow">
                <div class="grid grid-cols-[75%,20%,5%]">
                    <h5
                        class="flex text-lg text-left items-center font-bold text-gray-900"
                    >
                        <JiraIssueKey {issueKey} baseUrl={jiraBaseUrl} /> - {issue.summary}
                    </h5>
                    <StopWatch
                        time={timerValues[issueKey]}
                        startCallback={() => jira.StartTimer(issueKey)}
                        pauseCallback={() => jira.PauseTimer(issueKey)}
                        resetCallback={() => jira.ResetTimer(issueKey)}
                        setupCallback={() => setupTimerEventListener(issueKey)}
                    />
                    <div class="flex items-center">
                        {#if timerValues[issueKey] > 60 * 1e9}
                            <Button
                                size="xs"
                                class="h-10"
                                on:click={() =>
                                    jira
                                        .SubmitWorklog(issueKey)
                                        .catch((err) => wails.LogError(err))}
                                ><Icons.ShareAllSolid />
                            </Button>
                        {:else}
                            <Button disabled size="xs" class="h-10"
                                ><Icons.ShareAllSolid />
                            </Button>
                        {/if}
                    </div>
                </div>
                {#await jira.GetWorkLogs(issueKey)}
                    <div class="text-center p-10">
                        <Spinner size="10" />
                    </div>
                {:then workLogs}
                    <WorkLogList {workLogs} />
                {/await}
            </Card>
        {/each}
    {/await}
</div>
