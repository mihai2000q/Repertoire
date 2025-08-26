import { reduxRender, withToastify } from '../../../../test-utils.tsx'
import PerfectRehearsalMenuItem from './PerfectRehearsalMenuItem.tsx'
import { screen } from '@testing-library/react'
import { Menu } from '@mantine/core'
import { userEvent } from '@testing-library/user-event'
import { setupServer } from 'msw/node'
import { AddPerfectSongRehearsalRequest } from '../../../../types/requests/SongRequests.ts'
import { http, HttpResponse } from 'msw'
import { expect } from 'vitest'
import { AddPerfectRehearsalsToArtistsRequest } from '../../../../types/requests/ArtistRequests.ts'
import { AddPerfectRehearsalsToAlbumsRequest } from '../../../../types/requests/AlbumRequests.ts'
import { AddPerfectRehearsalsToPlaylistsRequest } from '../../../../types/requests/PlaylistRequests.ts'

describe('Perfect Rehearsal Menu Item', () => {
  const server = setupServer()

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', async () => {
    const user = userEvent.setup()

    reduxRender(
      withToastify(
        <Menu opened={true}>
          <Menu.Dropdown>
            <PerfectRehearsalMenuItem id={''} closeMenu={vi.fn()} type={'song'} />
          </Menu.Dropdown>
        </Menu>
      )
    )

    expect(screen.getByRole('menuitem', { name: /perfect rehearsal/i })).toBeInTheDocument()

    await user.click(screen.getByRole('menuitem', { name: /perfect rehearsal/i }))

    expect(await screen.findByRole('button', { name: 'cancel' })).toBeInTheDocument()
    expect(await screen.findByRole('button', { name: 'confirm' })).toBeInTheDocument()
  })

  describe('on confirmation', () => {
    it('should send artists request', async () => {
      const user = userEvent.setup()

      let capturedRequest: AddPerfectRehearsalsToArtistsRequest
      server.use(
        http.post('/artists/perfect-rehearsals', async (req) => {
          capturedRequest = (await req.request.json()) as AddPerfectRehearsalsToArtistsRequest
          return HttpResponse.json({ message: 'it worked' })
        })
      )

      const artistId = 'some-id'
      const closeMenu = vi.fn()

      reduxRender(
        withToastify(
          <Menu opened={true}>
            <Menu.Dropdown>
              <PerfectRehearsalMenuItem id={artistId} closeMenu={closeMenu} type={'artist'} />
            </Menu.Dropdown>
          </Menu>
        )
      )

      await user.click(screen.getByRole('menuitem', { name: /perfect rehearsal/i }))
      await user.click(await screen.findByRole('button', { name: 'confirm' }))

      expect(await screen.findByText(/perfect rehearsal added/i)).toBeInTheDocument()
      expect(capturedRequest).toStrictEqual({ ids: [artistId] })
      expect(closeMenu).toHaveBeenCalledOnce()
    })

    it('should send albums request', async () => {
      const user = userEvent.setup()

      let capturedRequest: AddPerfectRehearsalsToAlbumsRequest
      server.use(
        http.post('/albums/perfect-rehearsals', async (req) => {
          capturedRequest = (await req.request.json()) as AddPerfectRehearsalsToAlbumsRequest
          return HttpResponse.json({ message: 'it worked' })
        })
      )

      const albumId = 'some-id'
      const closeMenu = vi.fn()

      reduxRender(
        withToastify(
          <Menu opened={true}>
            <Menu.Dropdown>
              <PerfectRehearsalMenuItem id={albumId} closeMenu={closeMenu} type={'album'} />
            </Menu.Dropdown>
          </Menu>
        )
      )

      await user.click(screen.getByRole('menuitem', { name: /perfect rehearsal/i }))
      await user.click(await screen.findByRole('button', { name: 'confirm' }))

      expect(await screen.findByText(/perfect rehearsal added/i)).toBeInTheDocument()
      expect(capturedRequest).toStrictEqual({ ids: [albumId] })
      expect(closeMenu).toHaveBeenCalledOnce()
    })

    it('should send song request', async () => {
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
              <PerfectRehearsalMenuItem id={songId} closeMenu={closeMenu} type={'song'} />
            </Menu.Dropdown>
          </Menu>
        )
      )

      await user.click(screen.getByRole('menuitem', { name: /perfect rehearsal/i }))
      await user.click(await screen.findByRole('button', { name: 'confirm' }))

      expect(await screen.findByText(/perfect rehearsal added/i)).toBeInTheDocument()
      expect(capturedRequest).toStrictEqual({ id: songId })
      expect(closeMenu).toHaveBeenCalledOnce()
    })

    it('should send playlists request', async () => {
      const user = userEvent.setup()

      let capturedRequest: AddPerfectRehearsalsToPlaylistsRequest
      server.use(
        http.post('/playlists/perfect-rehearsals', async (req) => {
          capturedRequest = (await req.request.json()) as AddPerfectRehearsalsToPlaylistsRequest
          return HttpResponse.json({ message: 'it worked' })
        })
      )

      const playlistId = 'some-id'
      const closeMenu = vi.fn()

      reduxRender(
        withToastify(
          <Menu opened={true}>
            <Menu.Dropdown>
              <PerfectRehearsalMenuItem id={playlistId} closeMenu={closeMenu} type={'playlist'} />
            </Menu.Dropdown>
          </Menu>
        )
      )

      await user.click(screen.getByRole('menuitem', { name: /perfect rehearsal/i }))
      await user.click(await screen.findByRole('button', { name: 'confirm' }))

      expect(await screen.findByText(/perfect rehearsal added/i)).toBeInTheDocument()
      expect(capturedRequest).toStrictEqual({ ids: [playlistId] })
      expect(closeMenu).toHaveBeenCalledOnce()
    })
  })
})
