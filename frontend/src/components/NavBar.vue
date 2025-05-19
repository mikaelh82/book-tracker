<template>
    <nav class="bg-white dark:bg-gray-800 shadow-md p-4">
      <div class="container mx-auto flex justify-between items-center">
        <h1 class="text-xl font-bold text-gray-900 dark:text-white">Book Tracker</h1>
        <div class="space-x-4">
          <router-link
            to="/"
            class="text-gray-600 dark:text-gray-300 hover:text-blue-500 dark:hover:text-blue-400"
            active-class="text-blue-500 dark:text-blue-400"
          >
            Collection
          </router-link>
          <router-link
            to="/stats"
            class="text-gray-600 dark:text-gray-300 hover:text-blue-500 dark:hover:text-blue-400"
            active-class="text-blue-500 dark:text-blue-400"
          >
            Statistics
          </router-link>
        </div>
      </div>
    </nav>
  </template>
  
  <script setup lang="ts">
  </script>