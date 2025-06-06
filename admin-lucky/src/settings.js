export default {
  /**
   * Framework version number
   */
  version: '3.8.2',
  /**
   * page title
   */
  title: import.meta.env.VITE_APP_TITLE,
  /**
   * Sidebar theme dark theme theme-dark, light theme theme-light
   */
  sideTheme: 'theme-dark',
  /**
   * Frame theme color value
   */
  theme: '',
  //'#068CC5',
  /**
   * Whether the system layout configuration
   */
  showSettings: false,

  /**
   * Whether to display the top navigation
   */
  topNav: false,

  /**
   * Whether to display tagsView
   */
  tagsView: true,

  /**
   * Whether to fix the head
   */
  fixedHeader: false,

  /**
   * Whether to display the logo
   */
  sidebarLogo: true,

  /**
   * Whether to display dynamic title
   */
  dynamicTitle: false,

  /**
   * @type {string | array} 'production' | ['production', 'development']
   * @description Need show err logs component.
   * The default is only used in the production env
   * If you want to also use it in dev, you can pass ['production', 'development']
   */
  errorLog: 'production',
  /**
   * Copyright Information
   */
  copyright: 'Copyright ©2025',
  /**
   * Whether to show the bottom bar
   */
  showFooter: true,
  /**
   * Whether to display the watermark
   */
  showWatermark: false,
  /**
   * Watermark copywriting
   */
  watermarkText: 'tg',
  /**
   * Whether to show other logins
   */
  showOtherLogin: false,
  /**
   * default size
   */
  defaultSize: 'default'
}