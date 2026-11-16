<template>
    <div class="tooltip">
        <el-tooltip effect="dark" placement="top">
            <div slot="content">
                <p v-for="(item,index) in messages" :key="index">{{item}}</p>
            </div>
            <i class="el-icon-bangzhu1" style="font-size:16px;color:#828282;margin-left:8px;"></i>
        </el-tooltip>
    </div>
</template>
<script>
export default {
    name:'promptMessages',
    data(){
        return {};
    },
    props:{
        messages:{
            type:Array,
            default(){
                return [];
            }
        }
    }
}
</script>