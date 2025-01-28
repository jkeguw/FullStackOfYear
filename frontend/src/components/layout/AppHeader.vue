<template>
  <header class="bg-white shadow">
    <nav class="container mx-auto px-4 py-2">
      <div class="flex justify-between items-center">
        <router-link to="/" class="text-xl font-bold">CPC</router-link>
        <div class="flex gap-4">
          <router-link
            v-for="item in navItems"
            :key="item.path"
            :to="item.path"
            class="hover:text-blue-500"
          >
            {{ item.name }}
          </router-link>
        </div>
      </div>
    </nav>
  </header>
</template>

<script setup lang="ts">
const navItems = [
  { name: 'devices', path: '/devices' },
  { name: 'comments', path: '/reviews' },
  { name: 'tools', path: '/tools' }
]
</script>