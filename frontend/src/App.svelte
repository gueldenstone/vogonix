<script lang="ts">
  import {
    Navbar,
    NavBrand,
    NavHamburger,
    NavUl,
    Button,
  } from "flowbite-svelte";

  import IssueList from "@lib/IssueList.svelte";

  let issueListRef: IssueList;

  function handleRefresh() {
    if (issueListRef) {
      issueListRef.refresh();
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
      <Button pill size="sm" on:click={handleRefresh}>Refresh</Button>
    </div>
    <NavUl></NavUl>
  </Navbar>
  <IssueList bind:this={issueListRef} />
</main>

<style>
  @import "tailwindcss/base";
  @import "tailwindcss/components";
  @import "tailwindcss/utilities";
</style>
