import { reduxRender, withToastify } from '../../../../test-utils.tsx'
import PerfectRehearsalsMenuItem from './PerfectRehearsalsMenuItem.tsx'
import { screen } from '@testing-library/react'
import { Menu } from '@mantine/core'
import { userEvent } from '@testing-library/user-event'
import { setupServer } from 'msw/node'
import { AddPerfectSongRehearsalsRequest } from '../../../../types/requests/SongRequests.ts'
import { http, HttpResponse } from 'msw'
import { expect } from 'vitest'
import { AddPerfectRehearsalsToArtistsRequest } from '../../../../types/requests/ArtistRequests.ts'
import { AddPerfectRehearsalsToAlbumsRequest } from '../../../../types/requests/AlbumRequests.ts'
import { AddPerfectRehearsalsToPlaylistsRequest } from '../../../../types/requests/PlaylistRequests.ts'

describe('Perfect Rehearsals Menu Item', () => {
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
            <PerfectRehearsalsMenuItem ids={[]} type={'songs'} closeMenu={vi.fn()} />
          </Menu.Dropdown>
        </Menu>
      )
    )

    expect(screen.getByRole('menuitem', { name: /perfect rehearsals/i })).toBeInTheDocument()

    await user.click(screen.getByRole('menuitem', { name: /perfect rehearsals/i }))

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

      const artistIds = ['some-id-1', 'some-id-2']
      const closeMenu = vi.fn()

      reduxRender(
        withToastify(
          <Menu opened={true}>
            <Menu.Dropdown>
              <PerfectRehearsalsMenuItem ids={artistIds} closeMenu={closeMenu} type={'artists'} />
            </Menu.Dropdown>
          </Menu>
        )
      )

      await user.click(screen.getByRole('menuitem', { name: /perfect rehearsals/i }))
      await user.click(await screen.findByRole('button', { name: 'confirm' }))

      expect(await screen.findByText(/perfect rehearsals added/i)).toBeInTheDocument()
      expect(capturedRequest).toStrictEqual({ ids: artistIds })
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

      const albumIds = ['some-id-1', 'some-id-2']
      const closeMenu = vi.fn()

      reduxRender(
        withToastify(
          <Menu opened={true}>
            <Menu.Dropdown>
              <PerfectRehearsalsMenuItem ids={albumIds} closeMenu={closeMenu} type={'albums'} />
            </Menu.Dropdown>
          </Menu>
        )
      )

      await user.click(screen.getByRole('menuitem', { name: /perfect rehearsals/i }))
      await user.click(await screen.findByRole('button', { name: 'confirm' }))

      expect(await screen.findByText(/perfect rehearsals added/i)).toBeInTheDocument()
      expect(capturedRequest).toStrictEqual({ ids: albumIds })
      expect(closeMenu).toHaveBeenCalledOnce()
    })

    it('should send song request and call onSuccess', async () => {
      const user = userEvent.setup()

      let capturedRequest: AddPerfectSongRehearsalsRequest
      server.use(
        http.post('/songs/perfect-rehearsals', async (req) => {
          capturedRequest = (await req.request.json()) as AddPerfectSongRehearsalsRequest
          return HttpResponse.json({ message: 'it worked' })
        })
      )

      const songIds = ['some-id-1', 'some-id-2']
      const closeMenu = vi.fn()
      const onSuccess = vi.fn()

      reduxRender(
        withToastify(
          <Menu opened={true}>
            <Menu.Dropdown>
              <PerfectRehearsalsMenuItem
                ids={songIds}
                closeMenu={closeMenu}
                onSuccess={onSuccess}
                type={'songs'}
              />
            </Menu.Dropdown>
          </Menu>
        )
      )

      await user.click(screen.getByRole('menuitem', { name: /perfect rehearsals/i }))
      await user.click(await screen.findByRole('button', { name: 'confirm' }))

      expect(await screen.findByText(/perfect rehearsals added/i)).toBeInTheDocument()
      expect(capturedRequest).toStrictEqual({ ids: songIds })
      expect(closeMenu).toHaveBeenCalledOnce()
      expect(onSuccess).toHaveBeenCalledOnce()
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

      const playlistIds = ['some-id-1', 'some-id-2']
      const closeMenu = vi.fn()

      reduxRender(
        withToastify(
          <Menu opened={true}>
            <Menu.Dropdown>
              <PerfectRehearsalsMenuItem
                ids={playlistIds}
                closeMenu={closeMenu}
                type={'playlists'}
              />
            </Menu.Dropdown>
          </Menu>
        )
      )

      await user.click(screen.getByRole('menuitem', { name: /perfect rehearsals/i }))
      await user.click(await screen.findByRole('button', { name: 'confirm' }))

      expect(await screen.findByText(/perfect rehearsals added/i)).toBeInTheDocument()
      expect(capturedRequest).toStrictEqual({ ids: playlistIds })
      expect(closeMenu).toHaveBeenCalledOnce()
    })
  })
})
