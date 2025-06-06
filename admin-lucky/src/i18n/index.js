import { createI18n } from 'vue-i18n'
// import useAppStore from '@/store/modules/app'
import { listLangByLocale } from '@/api/system/commonLang'


import jsCookie from 'js-cookie'
const language = computed(() => {
  // return useAppStore().lang
  return jsCookie.get('lang') || 'en'

})

import zhCn from './lang/zh-cn.json'
import eng from './lang/en.json'
import zhTw from './lang/zh-tw.json'
import kh from './lang/kh.json'

import pageLoginCn from './pages/login/zh-cn.json'
import pageLoginEn from './pages/login/en.json'
import pageLoginKh from './pages/login/kh.json'
import pageLoginTw from './pages/login/zh-tw.json'

// menu page
import pagemenuCn from './pages/menu/zh-cn'
import pagemenuEn from './pages/menu/en'
import pagemenuKh from './pages/menu/kh'
import pagemenuTw from './pages/menu/zh-tw'



const i18n = createI18n({
  // Inject the $t function globally
  globalInjection: true,
  fallbackLocale: 'en',
  locale: `${language.value}`, // language selected by default
  legacy: false, // To use Composition API mode, you need to set it to false
  //missingWarn: false,
  messages: {
    en: {
      ...eng,
      ...pageLoginEn,
      ...pagemenuEn
    },
    km: {
      ...kh,
      ...pageLoginKh,
      ...pagemenuKh
    },
    'zh-cn': {
      ...zhCn,
      ...pageLoginCn,
      ...pagemenuCn
    },
    'zh-tw': {
      ...zhTw,
      ...pageLoginTw,
      ...pagemenuTw
    },

    //... Add additional language support here
  }
})

const loadLocale = () => {
  i18n.global.mergeLocaleMessage(language.value, {})
  console.log(language.value)
  //i18n.global.setLocaleMessage(language.value)
  /*
  listLangByLocale(language.value).then(res => {
    const { code, data } = res
    if (code == 200) {
      i18n.global.mergeLocaleMessage(language.value, data)
    }
  })
  */
}
loadLocale()
//i18n.global.setLocaleMessage(language.value)
export default i18n;
