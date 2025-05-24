import type { DefineComponent } from 'vue'

declare module '@vue/runtime-core' {
  export interface GlobalComponents {
    ElAside: DefineComponent<{
      width?: string | number
    }>
    ElContainer: DefineComponent<{}>
    ElHeader: DefineComponent<{
      height?: string | number
    }>
    ElMain: DefineComponent<{}>
    ElMenu: DefineComponent<{
      defaultActive?: string
      collapse?: boolean
      router?: boolean
      backgroundColor?: string
      textColor?: string
      activeTextColor?: string
    }>
    ElSubMenu: DefineComponent<{
      index: string
    }>
    ElMenuItem: DefineComponent<{
      index: string
    }>
    ElIcon: DefineComponent<{}>
    ElButton: DefineComponent<{
      type?: string
    }>
    ElDropdown: DefineComponent<{
      trigger?: string
    }>
    ElDropdownMenu: DefineComponent<{}>
    ElDropdownItem: DefineComponent<{
      divided?: boolean
    }>
  }
} 