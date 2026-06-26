import {
  emptyArtist,
  emptySong,
  emptySongArrangement,
  reduxRender,
  withToastify
} from '../../test-utils.tsx'
import CustomRehearsalsModal from './CustomRehearsalsModal.tsx'
import Song from '../../types/models/Song.ts'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import WithTotalCountResponse from '../../types/responses/WithTotalCountResponse.ts'
import { userEvent } from '@testing-library/user-event'
import { screen } from '@testing-library/react'
import { AddCustomSongRehearsalRequest } from '../../types/requests/SongRequests.ts'
import { expect } from 'vitest'

describe('Custom Rehearsals Modal', () => {
  const songs: Song[] = [
    {
      ...emptySong,
      id: '1',
      title: 'Empty Song',
      arrangements: []
    },
    {
      ...emptySong,
      id: '2',
      title: 'Song Without Default Arrangement',
      artist: {
        ...emptyArtist,
        id: '1',
        name: 'Artist Name'
      },
      imageUrl: 'something.png',
      arrangements: [
        {
          ...emptySongArrangement,
          id: '1',
          name: 'Arrangement 1'
        },
        {
          ...emptySongArrangement,
          id: '2',
          name: 'Arrangement 2'
        }
      ]
    },
    {
      ...emptySong,
      id: '3',
      title: 'Song With Default Arrangement',
      defaultArrangementId: '12',
      arrangements: [
        {
          ...emptySongArrangement,
          id: '12',
          name: 'Perfect Arrangement'
        }
      ]
    }
  ]
  const ids = songs.map((song) => song.id)

  const handlers = [
    http.get(`/songs`, () => {
      const res: WithTotalCountResponse<Song> = {
        models: songs,
        totalCount: songs.length
      }
      return HttpResponse.json(res)
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', async () => {
    const user = userEvent.setup()

    reduxRender(<CustomRehearsalsModal opened={true} onClose={vi.fn()} ids={ids} />)

    expect(await screen.findByRole('dialog', { name: 'Custom Rehearsals' })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: 'Custom Rehearsals' })).toBeInTheDocument()

    for (const song of songs) {
      if (song.imageUrl) {
        expect(screen.getByRole('img', { name: song.title })).toBeInTheDocument()
      } else {
        expect(screen.getByLabelText(`default-icon-${song.title}`)).toBeInTheDocument()
      }
      expect(screen.getByText(song.title)).toBeInTheDocument()
      if (song.artist) {
        expect(screen.getByText(song.artist.name)).toBeInTheDocument()
      }

      const arrangement =
        song.arrangements.length === 0
          ? null
          : song.defaultArrangementId
            ? song.arrangements.find((a) => a.id === song.defaultArrangementId)
            : song.arrangements[0]

      if (!arrangement) {
        expect(screen.getByRole('button', { name: /no arrangement found/i })).toBeInTheDocument()
        expect(screen.getByRole('button', { name: /no arrangement found/i })).toBeDisabled()
      } else {
        expect(screen.getByRole('button', { name: arrangement.name })).toBeInTheDocument()
        await user.click(screen.getByRole('button', { name: arrangement.name }))
        song.arrangements.forEach((a) => {
          expect(screen.getByRole('menuitem', { name: a.name })).toBeInTheDocument()
          if (a.id === arrangement.id) {
            expect(screen.getByRole('menuitem', { name: a.name })).toHaveAttribute(
              'data-active',
              'true'
            )
          }
        })
      }
    }

    expect(screen.getByRole('button', { name: /submit/i })).toBeInTheDocument()
  })

  it('should be able to switch song arrangements', async () => {
    const user = userEvent.setup()

    reduxRender(<CustomRehearsalsModal opened={true} onClose={vi.fn()} ids={ids} />)

    await user.click(await screen.findByRole('button', { name: songs[1].arrangements[0].name }))
    await user.click(await screen.findByRole('menuitem', { name: songs[1].arrangements[1].name }))
    expect(screen.getByRole('button', { name: songs[1].arrangements[1].name })).toBeInTheDocument()
    expect(
      screen.queryByRole('button', { name: songs[1].arrangements[0].name })
    ).not.toBeInTheDocument()
  })

  it('should be able to submit custom rehearsals', async () => {
    const user = userEvent.setup()

    let capturedRequest: AddCustomSongRehearsalRequest
    server.use(
      http.post('/songs/custom-rehearsals', async (req) => {
        capturedRequest = (await req.request.json()) as AddCustomSongRehearsalRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    const onClose = vi.fn()

    reduxRender(withToastify(<CustomRehearsalsModal opened={true} onClose={onClose} ids={ids} />))

    await user.click(await screen.findByRole('button', { name: /submit/i }))

    const requests: AddCustomSongRehearsalRequest[] = songs
      .filter((s) => s.arrangements.length > 0)
      .map((song) => ({
        id: song.id,
        arrangementId: song.defaultArrangementId
          ? song.arrangements.find((a) => a.id === song.defaultArrangementId).id
          : song.arrangements[0].id
      }))

    expect(await screen.findByText(/custom rehearsals added/i)).toBeInTheDocument()
    expect(capturedRequest).toStrictEqual({ requests: requests })
    expect(onClose).toHaveBeenCalledOnce()
  })
})
