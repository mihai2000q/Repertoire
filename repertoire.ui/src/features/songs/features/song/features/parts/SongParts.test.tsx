import { emptySongPart, reduxRender, withToastify } from '../../../../../../test-utils.tsx'
import SongParts from './SongParts.tsx'
import { SongPart } from '../../../../../../types/models/Song.ts'
import { fireEvent, screen, within } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { expect } from 'vitest'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import { MoveSongPartInSongRequest } from './types/requests/SongPartRequests.ts'
import { createRef } from 'react'

// Mock Main Context
vi.mock('../../../../../../context/MainContext.tsx', () => ({
  useMain: vi.fn(() => ({
    ref: createRef(),
    mainScroll: { ref: createRef() }
  }))
}))

describe('Song Parts', () => {
  const parts: SongPart[] = [
    {
      ...emptySongPart,
      id: '1',
      name: 'Chorus 1',
      rehearsals: 0,
      confidence: 0,
      progress: 0
    },
    {
      ...emptySongPart,
      id: '2',
      name: 'James Solo',
      rehearsals: 7,
      confidence: 50,
      progress: 163
    },
    {
      ...emptySongPart,
      id: '3',
      name: 'James Riff',
      rehearsals: 1,
      confidence: 36,
      progress: 40
    }
  ]

  const handlers = [
    http.get('/songs/instruments', () => {
      return HttpResponse.json([])
    }),
    http.put(`/songs/parts`, () => {
      return HttpResponse.json({ message: 'it worked' })
    })
  ]

  const server = setupServer(...handlers)

  afterEach(() => {
    server.resetHandlers()
    vi.restoreAllMocks()
  })

  beforeAll(() => server.listen())

  afterAll(() => server.close())

  it('should render', async () => {
    reduxRender(<SongParts parts={parts} />)

    const renderedParts = await screen.findAllByLabelText(/song-part-(?!details)/)
    for (let i = 0; i < parts.length; i++) {
      expect(renderedParts[i]).toHaveAccessibleName(`song-part-${parts[i].name}`)
    }
    screen.queryAllByLabelText(/song-part-details-/).forEach((d) => expect(d).not.toBeVisible())
  })

  it.skip('should be able to reorder parts', async () => {
    const part = parts[0]
    const overPart = parts[2]
    const songId = 'some-id'

    let capturedRequest: MoveSongPartInSongRequest
    server.use(
      http.put('/songs/sections/types', async (req) => {
        capturedRequest = (await req.request.json()) as MoveSongPartInSongRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    reduxRender(<SongParts parts={parts} />, {
      song: { songId: songId, isArtistBand: false }
    })

    fireEvent.mouseDown(screen.getByLabelText(`song-part-${part.name}`))
    fireEvent.dragStart(screen.getByLabelText(`song-part-${part.name}`))
    fireEvent.dragOver(screen.getByLabelText(`song-part-${overPart.name}`))
    fireEvent.drop(screen.getByLabelText(`song-part-${overPart.name}`))
    fireEvent.mouseUp(screen.getByLabelText(`song-part-${overPart.name}`))

    const renderedParts = await screen.findAllByLabelText(/song-part-(?!details)/)
    const expectedParts = [parts[1], parts[2], parts[0]]

    for (let i = 0; i < parts.length; i++) {
      expect(renderedParts[i]).toHaveAccessibleName(`song-part-${expectedParts[i].name}`)
    }

    expect(capturedRequest).toStrictEqual({
      id: part.id,
      overId: overPart.id,
      songId: songId
    })
  })

  it('should show toast when adding 1 rehearsal to part and dismiss on when adding to another part', async () => {
    const user = userEvent.setup()

    const part1 = parts[0]
    const part2 = parts[1]

    reduxRender(withToastify(<SongParts parts={parts} />))

    // click the first part
    await user.click(
      within(screen.getByLabelText(`song-part-${part1.name}`)).getByRole('button', {
        name: 'add-rehearsal'
      })
    )
    expect(
      await screen.findByText(new RegExp(`${part1.name} rehearsals.*increased.*1`, 'i'))
    ).toBeInTheDocument()

    // click the second part
    await user.click(
      within(screen.getByLabelText(`song-part-${part2.name}`)).getByRole('button', {
        name: 'add-rehearsal'
      })
    )

    expect(
      await screen.findByText(new RegExp(`${part2.name} rehearsals.*increased.*1`, 'i'))
    ).toBeInTheDocument()
    // not working - re-enable when switched to mantine notifs
    // expect(
    //   screen.queryByText(new RegExp(`${part1.name} rehearsals.*increased.*1`, 'i'))
    // ).not.toBeInTheDocument()
  })

  it('should show the drawer when selecting parts, and the context menu when right-clicking after selection', async () => {
    const user = userEvent.setup()

    reduxRender(<SongParts parts={parts} />)

    // selection drawer
    await user.keyboard('{Control>}')
    await user.click(screen.getByLabelText(`song-part-${parts[0].name}`))
    await user.keyboard('{/Control}')
    expect(screen.getByLabelText('song-parts-selection-drawer')).toBeInTheDocument()

    // context menu
    // await user.pointer({
    //   keys: '[MouseRight>]',
    //   target: screen.getByLabelText(`song-part-card-${songs[0].title}`)
    // })
    // expect(await screen.findByRole('menu', { name: 'song-parts-context-menu' })).toBeInTheDocument()
  })
})
