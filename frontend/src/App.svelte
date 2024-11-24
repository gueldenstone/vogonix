<script lang="ts">
  import "./app.css";
  // Importing the Flowbite Card component
  import {
    Navbar,
    NavBrand,
    NavHamburger,
    NavUl,
    NavLi,
    Button,
    Card,
    Spinner,
  } from "flowbite-svelte";

  import * as jira from "../wailsjs/go/jira/JiraInstance";
  import * as wails from "../wailsjs/runtime/runtime";
  import JiraIssueKey from "./lib/IssueKey.svelte";
  import StopWatch from "./lib/StopWatch.svelte";
  import * as Icons from "flowbite-svelte-icons";

  // Example array of issues
  let jiraBaseUrl;
  jira.GetBaseUrl().then((url) => (jiraBaseUrl = url));
  let issues = [];
  function updateIssues() {
    jira.GetAssignedIssues().then((jira_issues) => (issues = jira_issues));
  }
  // initialize
  updateIssues();

  let timerValues = new Map<string, number>();
  let stopWatchRefs: StopWatch[] = [];

  // Function to update the timer value when the stopwatch changes
  function handleStopwatchUpdate(issueKey: string, newValue: number) {
    wails.LogDebug(
      "Updating timer values for " + issueKey + " with " + newValue,
    );
    timerValues.set(issueKey, newValue);
  }

  // Submit the timer value for a specific issue
  function submitTime(issueKey: string) {
    const timeSpent = timerValues.get(issueKey);
    if (timeSpent !== undefined && timeSpent > 0) {
      jira.SubmitWorklfowLogs(issueKey, timeSpent);
    }
  }
</script>

<main>
  <Navbar>
    <NavBrand href="/">
      <span
        class="self-center whitespace-nowrap text-xl font-semibold dark:text-white"
        >vogonix</span
      >
    </NavBrand>
    <div class="flex md:order-2">
      <Button pill size="sm" on:click={updateIssues}>Refresh</Button>
      <NavHamburger />
    </div>
    <NavUl></NavUl>
  </Navbar>
  <div class="flex p-2 gap-2 flex-col">
    {#each issues as issue, index (issue.Key)}
      <Card size="xl" class="h-full grow">
        <div class="grid grid-cols-[70%,5%,20%,5%]">
          <h5
            class="flex text-lg text-left items-center font-bold text-gray-900"
          >
            <JiraIssueKey issueKey={issue.Key} baseUrl={jiraBaseUrl} /> - {issue.Summary}
          </h5>
          <p class="text-sm text-gray-700 flex items-center text-left">
            {#await jira.GetTimeSpentOnIssue(issue.Key)}
              <Spinner />
            {:then time}
              {time}
            {:catch error}
              Error!
            {/await}
          </p>
          <StopWatch
            timer={timerValues.get(issue.Key)}
            bind:this={stopWatchRefs[index]}
            updateCallback={() => handleStopwatchUpdate(index)}
          />
          <div class="flex items-center">
            <Button
              size="xs"
              class="h-10"
              on:click={() => {
                submitTime(issue.Key);
                stopWatchRefs[index].resetTimer();
              }}><Icons.ShareAllSolid /></Button
            >
          </div>
        </div>
      </Card>
    {/each}
  </div>
</main>

<style>
  /* Adjust the grid layout for responsive display */
</style>
