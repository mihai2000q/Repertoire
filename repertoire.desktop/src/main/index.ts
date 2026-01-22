import { app, BrowserWindow, ipcMain, shell } from 'electron'
import { join } from 'path'
import { electronApp, is, optimizer } from '@electron-toolkit/utils'

function createWindow(): void {
  // Splash Screen/Window
  const splashStartTime = Date.now()
  const splash = new BrowserWindow({
    center: true,
    transparent: true,
    frame: false,
    alwaysOnTop: true,
    titleBarStyle: 'hidden',
    autoHideMenuBar: true,
    webPreferences: { sandbox: false }
  })
  splash.removeMenu()

  // Create the main browser window.
  const mainWindow = new BrowserWindow({
    width: 1120,
    height: 800,
    minWidth: 460,
    minHeight: 460,
    show: false,
    autoHideMenuBar: true,
    webPreferences: {
      preload: join(__dirname, '../preload/index.js'),
      sandbox: false
    },
    titleBarStyle: 'hidden'
  })
  mainWindow.removeMenu()
  mainWindow.webContents.openDevTools()

  // HMR for renderer base on electron-vite cli.
  // Load the remote URL for development or the local HTML files for production.
  if (is.dev && process.env['ELECTRON_RENDERER_URL']) {
    mainWindow.loadURL(`${process.env['ELECTRON_RENDERER_URL']}/index.html`).then()
    splash.loadURL(`${process.env['ELECTRON_RENDERER_URL']}/splash.html`).then()
  } else {
    mainWindow.loadFile(join(__dirname, '../renderer/index.html')).then()
    splash.loadFile(join(__dirname, '../renderer/splash.html')).then()
  }

  mainWindow.on('ready-to-show', () => {
    const splashElapsed = Date.now() - splashStartTime
    const splashRemainingDelay = Math.max(
      0,
      parseInt(import.meta.env['SPLASH_SCREEN_MINIMUM_TIME']) - splashElapsed
    )

    setTimeout(() => {
      splash.destroy()
      mainWindow.show()
    }, splashRemainingDelay)
  })

  mainWindow.webContents.setWindowOpenHandler((details) => {
    shell.openExternal(details.url).then()
    return { action: 'deny' }
  })

  // set CSP
  mainWindow.webContents.session.webRequest.onHeadersReceived((details, callback) => {
    callback({
      responseHeaders: {
        ...details.responseHeaders,
        'Content-Security-Policy': [
          // TODO: Replace youtube-modal with proper encrypted nonce
          `default-src 'self';
          script-src 'self' https://*.google.com https://*.gstatic.com 'unsafe-inline' 'unsafe-eval';
          style-src 'self' ${is.dev ? 'unsafe-inline' : 'nonce-youtube-modal'} 'unsafe-inline' https://fonts.googleapis.com;
          font-src 'self' https://fonts.gstatic.com;
          img-src 'self' ${import.meta.env.VITE_WEB_ORIGINS} https: data: blob:;
          connect-src 'self' https://*.googleapis.com https://*.googlevideo.com ${import.meta.env.VITE_WEB_ORIGINS};
          frame-src 'self' https://*.youtube-nocookie.com;
          media-src 'self' https://*.youtube-nocookie.com blob:;
          object-src 'none';
          child-src 'self' https://*.youtube-nocookie.com;`
        ],
        'Permissions-Policy': [
          'accelerometer=*, autoplay=*, camera=*, display-capture=*, encrypted-media=*, fullscreen=*, geolocation=*, gyroscope=*, keyboard-map=*, magnetometer=*, microphone=*, midi=*, payment=*, picture-in-picture=*, publickey-credentials-get=*, screen-wake-lock=*, sync-xhr=*, usb=*, xr-spatial-tracking=*'
        ]
      }
    })
  })

  mainWindow.webContents.setUserAgent(
    'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36'
  )

  // only 1 instance active at a time
  const gotTheLock = app.requestSingleInstanceLock()

  if (!gotTheLock) {
    app.quit()
  } else {
    app.on('second-instance', () => {
      // Someone tried to run a second instance, we should focus our window.
      if (mainWindow) {
        if (mainWindow.isMinimized()) mainWindow.restore()
        mainWindow.focus()
      }
    })
  }

  // ipc handlers
  ipcMain.on('minimize', () =>
    mainWindow.isMinimized() ? mainWindow.restore() : mainWindow.minimize()
  )
  ipcMain.on('maximize', () =>
    mainWindow.isMaximized() ? mainWindow.unmaximize() : mainWindow.maximize()
  )
  ipcMain.on('close', () => mainWindow.close())
}

// This method will be called when Electron has finished
// initialization and is ready to create browser windows.
// Some APIs can only be used after this event occurs.
app.whenReady().then(() => {
  // Set app user model id for windows
  electronApp.setAppUserModelId('com.electron.repertoire')

  // Default open or close DevTools by F12 in development
  // and ignore CommandOrControl + R in production.
  // See https://github.com/alex8088/electron-toolkit/tree/master/packages/utils
  app.on('browser-window-created', (_, window) => {
    optimizer.watchWindowShortcuts(window)
  })

  // IPC
  ipcMain.on('log', (_, log) => console.log(log))

  createWindow()

  app.on('activate', function () {
    // On macOS, it's common to re-create a window in the app when the
    // dock icon is clicked and there are no other windows open.
    if (BrowserWindow.getAllWindows().length === 0) createWindow()
  })
})

// Quit when all windows are closed, except on macOS. There, it's common
// for applications and their menu bar to stay active until the user quits
// explicitly with Cmd + Q.
app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit()
  }
})
