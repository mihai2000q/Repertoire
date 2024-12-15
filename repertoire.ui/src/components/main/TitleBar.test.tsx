import { mantineRender } from '../../test-utils.tsx'
import TitleBar from './TitleBar.tsx'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { afterEach } from 'vitest'
import { ElectronAPI, IpcRenderer } from '@electron-toolkit/preload'

describe('Title Bar', () => {
  afterEach(() => vi.restoreAllMocks())

  it('should render and display 3 action buttons', () => {
    mantineRender(<TitleBar />)

    expect(screen.getByRole('button', { name: 'minimize' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'maximize' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'close' })).toBeInTheDocument()
  })

  it('should have 3 action buttons that trigger the electron Api', async () => {
    // Arrange
    const user = userEvent.setup()
    window.electron = {
      ipcRenderer: {
        send: vi.fn()
      } as unknown as IpcRenderer
    } as unknown as ElectronAPI

    // Act
    mantineRender(<TitleBar />)

    // Act & Assert
    await user.click(screen.getByRole('button', { name: 'minimize' }))
    expect(window.electron.ipcRenderer.send).toHaveBeenCalledWith('minimize')

    await user.click(screen.getByRole('button', { name: 'maximize' }))
    expect(window.electron.ipcRenderer.send).toHaveBeenCalledWith('maximize')

    await user.click(screen.getByRole('button', { name: 'close' }))
    expect(window.electron.ipcRenderer.send).toHaveBeenCalledWith('close')
    expect(window.electron.ipcRenderer.send).toHaveBeenCalledTimes(3)
  })
})
