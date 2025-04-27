<template>
  <div class="contact-page container mx-auto py-8 px-4">
    <h1 class="text-3xl font-bold mb-8 text-white">联系我们</h1>
    
    <div class="grid md:grid-cols-2 gap-8">
      <div class="bg-[#1E1E1E] rounded-lg p-6 text-gray-300">
        <h2 class="text-2xl font-bold mb-4 text-white">联系方式</h2>
        
        <div class="mb-4">
          <h3 class="text-xl mb-2 text-white">客户支持</h3>
          <p>邮箱：support@mousecomparison.com</p>
          <p>工作时间：周一至周五 9:00-18:00</p>
        </div>
        
        <div class="mb-4">
          <h3 class="text-xl mb-2 text-white">商务合作</h3>
          <p>邮箱：business@mousecomparison.com</p>
        </div>
        
        <div>
          <h3 class="text-xl mb-2 text-white">社交媒体</h3>
          <div class="flex space-x-4 mt-2">
            <a href="#" class="text-blue-400 hover:text-blue-300">微博</a>
            <a href="#" class="text-blue-400 hover:text-blue-300">知乎</a>
            <a href="#" class="text-blue-400 hover:text-blue-300">B站</a>
            <a href="#" class="text-blue-400 hover:text-blue-300">Twitter</a>
          </div>
        </div>
      </div>
      
      <div class="bg-[#1E1E1E] rounded-lg p-6 text-gray-300">
        <h2 class="text-2xl font-bold mb-4 text-white">给我们留言</h2>
        
        <form class="space-y-4">
          <div>
            <label class="block mb-2 text-white">您的姓名</label>
            <el-input placeholder="请输入您的姓名" />
          </div>
          
          <div>
            <label class="block mb-2 text-white">电子邮箱</label>
            <el-input placeholder="请输入您的电子邮箱" />
          </div>
          
          <div>
            <label class="block mb-2 text-white">留言内容</label>
            <el-input type="textarea" :rows="5" placeholder="请输入您的留言内容" />
          </div>
          
          <el-button type="primary" class="w-full">提交</el-button>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// 无需额外脚本
</script>

<style scoped>
/* 使用全局变量确保与主题一致 */
</style>