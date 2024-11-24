<script lang="ts">
    import { Button, Card, Spinner } from "flowbite-svelte";
    import * as jira from "@go/jira/JiraInstance";
    import * as wails from "@runtime/runtime";

    import JiraIssueKey from "@lib/IssueKey.svelte";
    import StopWatch from "@lib/StopWatch.svelte";
    import * as Icons from "flowbite-svelte-icons";
    import { onMount } from "svelte";

    // base url
    let jiraBaseUrl: string;
    jira.GetBaseUrl().then((url) => (jiraBaseUrl = url));

    // local state for timer values
    $: timerValues = new Map<string, number>();

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
            timerValues[issue.Key] = await jira.GetCurrentTimerValue(issue.Key);
        });
    });
</script>

<div class="flex p-2 gap-2 flex-col">
    {#await issuesPromise}
        <div class="text-center p-10">
            <Spinner size="10" />
        </div>
    {:then issues}
        {#each issues as issue, index (issue.Key)}
            {@const issueKey = issue.Key}
            <Card size="xl" class="h-full grow">
                <div class="grid grid-cols-[70%,5%,20%,5%]">
                    <h5
                        class="flex text-lg text-left items-center font-bold text-gray-900"
                    >
                        <JiraIssueKey {issueKey} baseUrl={jiraBaseUrl} /> - {issue.Summary}
                    </h5>
                    <p
                        class="text-sm text-gray-700 flex items-center text-left"
                    >
                        {#await jira.GetTimeSpentOnIssue(issueKey)}
                            <Spinner />
                        {:then time}
                            {time}
                        {:catch error}
                            Error!
                        {/await}
                    </p>
                    <StopWatch
                        time={timerValues[issueKey]}
                        startCallback={() => jira.StartTimer(issueKey)}
                        pauseCallback={() => jira.PauseTimer(issueKey)}
                        resetCallback={() => jira.ResetTimer(issueKey)}
                        setupCallback={() => setupTimerEventListener(issueKey)}
                    />
                    <div class="flex items-center">
                        <Button
                            size="xs"
                            class="h-10"
                            on:click={() => jira.SubmitWorklog(issueKey)}
                            ><Icons.ShareAllSolid /></Button
                        >
                    </div>
                </div>
            </Card>
        {/each}
    {/await}
</div>
