<script lang="ts">
    import { Button, Card, Spinner } from "flowbite-svelte";
    import JiraIssueKey from "@lib/IssueKey.svelte";
    import StopWatch from "@lib/StopWatch.svelte";
    import * as Icons from "flowbite-svelte-icons";
    import WorkLogList from "./WorkLogList.svelte";
    import * as jira from "@go/jira/JiraInstance";
    import * as wails from "@runtime/runtime";
    import * as models from "@go/models";

    export let timerValue;
    export let issue: models.jira.Issue;

    $: issueKey = issue.key;

    // base url
    let jiraBaseUrl: string;
    jira.GetBaseUrl().then((url) => (jiraBaseUrl = url));

    function setupTimerEventListener() {
        wails.EventsOn("timer_tick_" + issueKey, (currentTime) => {
            timerValue = currentTime;
        });
    }
</script>

<Card size="xl" class="h-full grow">
    <div class="grid grid-cols-[75%,20%,5%]">
        <h5 class="flex text-lg text-left items-center font-bold text-gray-900">
            <JiraIssueKey {issueKey} baseUrl={jiraBaseUrl} /> - {issue.summary}
        </h5>
        <StopWatch
            time={timerValue}
            startCallback={() => jira.StartTimer(issueKey)}
            pauseCallback={() => jira.PauseTimer(issueKey)}
            resetCallback={() => jira.ResetTimer(issueKey)}
            setupCallback={setupTimerEventListener}
        />
        <div class="flex items-center">
            {#if timerValue > 60 * 1e9}
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
