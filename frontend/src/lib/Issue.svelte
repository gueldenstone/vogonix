<script lang="ts">
    import { Card } from "flowbite-svelte";
    import JiraIssueKey from "@lib/IssueKey.svelte";
    import StopWatch from "@lib/StopWatch.svelte";
    import WorkLogList from "@lib/WorkLogList.svelte";
    import * as jira from "@go/jira/JiraInstance";
    import * as wails from "@runtime/runtime";
    import * as models from "@go/models";
    import SubmitButton from "@lib/SubmitButton.svelte";

    export let timerValue: number;
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

    $: submitEnabled = timerValue > 60 * 1e9;
</script>

<Card size="xl" class="grow">
    <div class="grid grid-cols-[60%,30%,10%] gap-1">
        <h5 class="flex text-lg text-left items-center font-bold text-gray-900">
            <JiraIssueKey {issueKey} baseUrl={jiraBaseUrl} /> -
            <span class="truncate"> {issue.summary}</span>
        </h5>
        <StopWatch
            time={timerValue}
            startCallback={() => jira.StartTimer(issueKey)}
            pauseCallback={() => jira.PauseTimer(issueKey)}
            resetCallback={() => jira.ResetTimer(issueKey)}
            setupCallback={setupTimerEventListener}
        />
        <div class="flex items-center">
            <SubmitButton enabled={submitEnabled} {issueKey} />
        </div>
    </div>
    <WorkLogList {issueKey} />
</Card>
