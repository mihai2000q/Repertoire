import { reduxRender, withToastify } from '../../../../../test-utils.tsx'
import PartialRehearsalMenuItem from './PartialRehearsalMenuItem.tsx'
import { screen } from '@testing-library/react'
import { Menu } from '@mantine/core'
import { userEvent } from '@testing-library/user-event'
import { setupServer } from 'msw/node'
import { AddPartialSongRehearsalRequest } from '../../../../../types/requests/SongRequests.ts'
import { http, HttpResponse } from 'msw'
import { expect } from 'vitest'

describe('Partial Rehearsal Menu Item', () => {
  const server = setupServer()

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render and send request on confirmation', async () => {
    const user = userEvent.setup()

    let capturedRequest: AddPartialSongRehearsalRequest
    server.use(
      http.post('/songs/partial-rehearsal', async (req) => {
        capturedRequest = (await req.request.json()) as AddPartialSongRehearsalRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    const songId = 'some-id'

    reduxRender(
      withToastify(
        <Menu opened={true}>
          <Menu.Dropdown>
            <PartialRehearsalMenuItem songId={songId} />
          </Menu.Dropdown>
        </Menu>
      )
    )

    await user.click(screen.getByRole('menuitem', { name: /partial rehearsal/i }))
    await user.click(screen.getByRole('button', { name: 'confirm' }))

    expect(await screen.findByText(/partial rehearsal added/i)).toBeInTheDocument()
    expect(capturedRequest).toStrictEqual({ id: songId })
  })
})
