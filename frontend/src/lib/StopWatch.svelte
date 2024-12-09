<script lang="ts">
    import { Button, Badge } from "flowbite-svelte";
    import * as Icons from "flowbite-svelte-icons";
    import { onMount } from "svelte";
    import { formatDurationFromNanoseconds } from "@lib/util";

    // export let time: number = 0;
    export let isRunning: boolean = false;

    export let startCallback;
    export let pauseCallback;
    export let resetCallback;
    export let setupCallback;
    // export let newTimeCallback;

    export let time: number = 0;

    onMount(() => {
        setupCallback();
    });
</script>

<div class="grid grid-cols-[30%,70%] gap-2 align-center items-center">
    <!-- Display the timer -->
    <Badge color="dark" border class=" min-h-8">
        <!-- <Icons.ClockSolid class="w-2.5 h-2.5 me-1.5" /> -->
        {formatDurationFromNanoseconds(time)}
    </Badge>

    <!-- Control buttons -->
    <div class="grid grid-cols-2 gap-2 m-4">
        {#if !isRunning}
            <Button
                on:click={() => {
                    startCallback();
                    isRunning = true;
                }}
                color="light"
                size="xs"
            >
                <Icons.PlaySolid />
            </Button>
        {:else}
            <Button
                on:click={() => {
                    pauseCallback();
                    isRunning = false;
                }}
                color="light"
                size="xs"
            >
                <Icons.PauseSolid />
            </Button>
        {/if}
        <Button
            on:click={() => {
                resetCallback();
                isRunning = false;
            }}
            color="light"
            size="xs"
        >
            <Icons.CloseCircleSolid />
        </Button>
    </div>
</div>
