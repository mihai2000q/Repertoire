import { reduxRender, withToastify } from '../../../../../test-utils.tsx'
import PerfectRehearsalMenuItem from './PerfectRehearsalMenuItem.tsx'
import { screen } from '@testing-library/react'
import { Menu } from '@mantine/core'
import { userEvent } from '@testing-library/user-event'
import { setupServer } from 'msw/node'
import { AddPerfectSongRehearsalRequest } from '../../../../../types/requests/SongRequests.ts'
import { http, HttpResponse } from 'msw'
import { expect } from 'vitest'

describe('Perfect Rehearsal Menu Item', () => {
  const server = setupServer()

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render and send request on confirmation', async () => {
    const user = userEvent.setup()

    let capturedRequest: AddPerfectSongRehearsalRequest
    server.use(
      http.post('/songs/perfect-rehearsal', async (req) => {
        capturedRequest = (await req.request.json()) as AddPerfectSongRehearsalRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    const songId = 'some-id'
    const closeMenu = vi.fn()

    reduxRender(
      withToastify(
        <Menu opened={true}>
          <Menu.Dropdown>
            <PerfectRehearsalMenuItem songId={songId} closeMenu={closeMenu} />
          </Menu.Dropdown>
        </Menu>
      )
    )

    await user.click(screen.getByRole('menuitem', { name: /perfect rehearsal/i }))
    await user.click(screen.getByRole('button', { name: 'confirm' }))

    expect(await screen.findByText(/perfect rehearsal added/i)).toBeInTheDocument()
    expect(capturedRequest).toStrictEqual({ id: songId })
    expect(closeMenu).toHaveBeenCalledOnce()
  })
})
