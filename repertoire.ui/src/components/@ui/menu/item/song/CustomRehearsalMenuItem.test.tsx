import { userEvent } from '@testing-library/user-event/dist/cjs/index.js'
import { emptySongArrangement, reduxRender, withToastify } from '../../../../../test-utils.tsx'
import { Menu } from '@mantine/core'
import { expect } from 'vitest'
import { screen } from '@testing-library/react'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import CustomRehearsalMenuItem from './CustomRehearsalMenuItem.tsx'
import { SongArrangement } from '../../../../../types/models/Song.ts'
import { AddCustomSongRehearsalRequest } from '../../../../../types/requests/SongRequests.ts'

describe('Custom Rehearsal Menu Item', () => {
  const arrangements: SongArrangement[] = [
    {
      ...emptySongArrangement,
      id: '1',
      name: 'something'
    },
    {
      ...emptySongArrangement,
      id: '2',
      name: 'something 2'
    }
  ]

  const handlers = [
    http.get(`/songs/arrangements`, () => {
      return HttpResponse.json(arrangements)
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', async () => {
    const user = userEvent.setup()

    reduxRender(
      withToastify(
        <Menu opened={true}>
          <Menu.Dropdown>
            <CustomRehearsalMenuItem id={''} closeMenu={vi.fn()} />
          </Menu.Dropdown>
        </Menu>
      )
    )

    expect(screen.getByRole('menuitem', { name: /custom rehearsal/i })).toBeInTheDocument()
    await user.hover(screen.getByRole('menuitem', { name: /custom rehearsal/i }))

    for (const arrangement of arrangements) {
      expect(await screen.findByRole('menuitem', { name: arrangement.name })).toBeInTheDocument()
    }

    await user.click(screen.getByRole('menuitem', { name: arrangements[0].name }))

    expect(await screen.findByRole('button', { name: 'cancel' })).toBeInTheDocument()
    expect(await screen.findByRole('button', { name: 'confirm' })).toBeInTheDocument()
  })

  it('should display message when there are no arrangements', async () => {
    const user = userEvent.setup()

    server.use(
      http.get(`/songs/arrangements`, () => {
        return HttpResponse.json([])
      })
    )

    reduxRender(
      withToastify(
        <Menu opened={true}>
          <Menu.Dropdown>
            <CustomRehearsalMenuItem id={''} closeMenu={vi.fn()} />
          </Menu.Dropdown>
        </Menu>
      )
    )

    await user.hover(screen.getByRole('menuitem', { name: /custom rehearsal/i }))

    expect(await screen.findByText(/no arrangements found/i)).toBeInTheDocument()
  })

  it('should send custom rehearsal request when confirming on an arrangement', async () => {
    const user = userEvent.setup()

    let capturedRequest: AddCustomSongRehearsalRequest
    server.use(
      http.post('/songs/custom-rehearsal', async (req) => {
        capturedRequest = (await req.request.json()) as AddCustomSongRehearsalRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    const arrangement = arrangements[0]
    const songId = 'some-id'
    const closeMenu = vi.fn()

    reduxRender(
      withToastify(
        <Menu opened={true}>
          <Menu.Dropdown>
            <CustomRehearsalMenuItem id={songId} closeMenu={closeMenu} />
          </Menu.Dropdown>
        </Menu>
      )
    )

    await user.hover(screen.getByRole('menuitem', { name: /custom rehearsal/i }))

    await user.click(await screen.findByRole('menuitem', { name: arrangement.name }))
    await user.click(await screen.findByRole('button', { name: 'confirm' }))

    expect(await screen.findByText(/custom rehearsal added/i)).toBeInTheDocument()
    expect(capturedRequest).toStrictEqual({ id: songId, arrangementId: arrangement.id })
    expect(closeMenu).toHaveBeenCalledOnce()
  })
})
