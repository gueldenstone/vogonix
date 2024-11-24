<script lang="ts">
    import { Button } from "flowbite-svelte";
    import * as Icons from "flowbite-svelte-icons";
    import { onMount } from "svelte";

    export let timer: number = 0; // Timer value in seconds
    let interval: NodeJS.Timer | null = null;

    // Compute hours, minutes, and seconds
    $: hours = Math.floor(timer / 3600);
    $: minutes = Math.floor((timer % 3600) / 60);
    $: seconds = timer % 60;

    // Start the stopwatch
    const start = () => {
        if (!interval) {
            interval = setInterval(() => {
                timer++; // Increment the timer
            }, 1000);
        }
    };

    // Pause the stopwatch
    const pause = () => {
        console.log("pause");
        if (interval) {
            clearInterval(interval);
            interval = null;
        }
    };

    // Reset the stopwatch
    const reset = () => {
        pause();
        timer = 0;
    };

    // Allow parent component to reset the stopwatch
    export const resetTimer = () => {
        reset(); // Call internal reset method
    };

    export let updateCallback : 

    // Clean up the interval on component destroy
    onMount(() => {
        return () => {
            if (interval) clearInterval(interval);
        };
    });
</script>

<div class="stopwatch">
    <!-- Display the timer -->
    <div class="time-display">
        {hours.toString().padStart(2, "0")}:{minutes
            .toString()
            .padStart(2, "0")}:
        {seconds.toString().padStart(2, "0")}
    </div>

    <!-- Control buttons -->
    <div class="buttons">
        {#if !interval}
            <Button on:click={start} color="light" size="xs">
                <Icons.PlaySolid />
            </Button>
        {:else}
            <Button on:click={pause} color="light" size="xs">
                <Icons.PauseSolid />
            </Button>
        {/if}
        <Button on:click={reset} color="light" size="xs">
            <Icons.CloseCircleSolid />
        </Button>
    </div>
</div>
