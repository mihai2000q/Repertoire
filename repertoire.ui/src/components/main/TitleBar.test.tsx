import { reduxRender } from '../../test-utils.tsx'
import TitleBar from './TitleBar.tsx'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { afterEach } from 'vitest'
import { ElectronAPI, IpcRenderer } from '@electron-toolkit/preload'

describe('Title Bar', () => {
  afterEach(() => vi.restoreAllMocks())

  it('should render and display the document title and 3 action buttons', () => {
    const documentTitle = "Title bar's document title"

    reduxRender(<TitleBar />, {
      global: {
        documentTitle: documentTitle,
        artistDrawer: undefined,
        songDrawer: undefined,
        albumDrawer: undefined
      }
    })

    expect(window.document.title).toBe(documentTitle)
    expect(screen.getByText(documentTitle)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'minimize' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'maximize' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'close' })).toBeInTheDocument()
  })

  it('should have 3 action buttons that trigger the electron Api', async () => {
    const user = userEvent.setup()
    window.electron = {
      ipcRenderer: {
        send: vi.fn()
      } as unknown as IpcRenderer
    } as unknown as ElectronAPI

    reduxRender(<TitleBar />)

    await user.click(screen.getByRole('button', { name: 'minimize' }))
    expect(window.electron.ipcRenderer.send).toHaveBeenCalledWith('minimize')

    await user.click(screen.getByRole('button', { name: 'maximize' }))
    expect(window.electron.ipcRenderer.send).toHaveBeenCalledWith('maximize')

    await user.click(screen.getByRole('button', { name: 'close' }))
    expect(window.electron.ipcRenderer.send).toHaveBeenCalledWith('close')
    expect(window.electron.ipcRenderer.send).toHaveBeenCalledTimes(3)
  })
})
