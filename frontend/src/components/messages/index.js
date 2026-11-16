import promptMessages from './promptMessages'

promptMessages.install = function(Vue) {
    Vue.component(promptMessages.name,promptMessages);
}
export default promptMessages;