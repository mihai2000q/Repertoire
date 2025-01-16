import { ElectronAPI } from '@electron-toolkit/preload'

declare global {
  // noinspection JSUnusedGlobalSymbols
  interface Window {
    electron: ElectronAPI
    api: unknown
  }
}
