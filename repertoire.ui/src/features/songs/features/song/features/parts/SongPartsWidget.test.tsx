import {
  emptySongArrangement,
  emptySongPart,
  emptySongSettings,
  reduxRender,
  withToastify
} from '../../../../../../test-utils.tsx'
import SongPartsWidget from './SongPartsWidget.tsx'
import { SongArrangement, SongPart } from '../../../../../../types/models/Song.ts'
import { fireEvent, screen, within } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { expect } from 'vitest'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import { AddPerfectSongRehearsalRequest } from '../../../../../../types/requests/SongRequests.ts'
import { MoveSongPartInSongRequest } from './types/requests/SongPartRequests.ts'
import { createRef } from 'react'

// Mock Main Context
vi.mock('../../../../../../context/MainContext.tsx', () => ({
  useMain: vi.fn(() => ({
    ref: createRef(),
    mainScroll: { ref: createRef() }
  }))
}))

describe('Song Parts Widget', () => {
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

  const arrangements: SongArrangement[] = [
    {
      ...emptySongArrangement,
      id: '1'
    }
  ]

  const handlers = [
    http.get('/songs/sections/types', () => {
      return HttpResponse.json([])
    }),
    http.get('/songs/instruments', () => {
      return HttpResponse.json([])
    }),
    http.get(`/songs/arrangements`, () => {
      return HttpResponse.json(arrangements)
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
    reduxRender(
      <SongPartsWidget
        parts={parts}
        songId={''}
        settings={emptySongSettings}
        defaultSongArrangementId={'id'}
      />
    )

    expect(screen.getByText(/parts/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'add-new-part' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'show-details' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'show-details' })).not.toBeDisabled()
    expect(screen.getByRole('button', { name: 'manage-song-arrangements' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'add-custom-rehearsal' })).toBeInTheDocument()
    expect(await screen.findByRole('button', { name: 'add-custom-rehearsal' })).not.toBeDisabled()
    expect(screen.getByRole('button', { name: 'add-perfect-rehearsal' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'add-perfect-rehearsal' })).not.toBeDisabled()
    expect(screen.getByRole('button', { name: 'settings' })).toBeInTheDocument()

    const renderedParts = screen.getAllByLabelText(/song-part-(?!details)/)
    for (let i = 0; i < parts.length; i++) {
      expect(renderedParts[i]).toHaveAccessibleName(`song-part-${parts[i].name}`)
    }
    screen.queryAllByLabelText(/song-part-details-/).forEach((d) => expect(d).not.toBeVisible())
  })

  it('should disable a few options when there are no parts', () => {
    reduxRender(
      <SongPartsWidget
        parts={[]}
        songId={''}
        settings={emptySongSettings}
        defaultSongArrangementId={'id'}
      />
    )

    expect(screen.getByRole('button', { name: 'show-details' })).toBeDisabled()
    expect(screen.getByRole('button', { name: 'add-custom-rehearsal' })).toBeDisabled()
    expect(screen.getByRole('button', { name: 'add-perfect-rehearsal' })).toBeDisabled()
  })

  it('should disable perfect rehearsal when a default arrangement is not set', () => {
    reduxRender(
      <SongPartsWidget
        parts={parts}
        songId={''}
        settings={emptySongSettings}
        defaultSongArrangementId={null}
      />
    )

    expect(screen.getByRole('button', { name: 'add-perfect-rehearsal' })).toBeDisabled()
  })

  describe('on toolbar options', () => {
    it('should open add new song part when clicking on add new part button', async () => {
      const user = userEvent.setup()

      reduxRender(
        <SongPartsWidget parts={parts} songId={''} settings={emptySongSettings} />
      )

      await user.click(screen.getByRole('button', { name: 'add-new-part' }))
      expect(screen.getByLabelText('add-new-song-part')).toBeInTheDocument()
    })

    it('should show details when clicking on show details', async () => {
      const user = userEvent.setup()

      reduxRender(
        <SongPartsWidget parts={parts} songId={''} settings={emptySongSettings} />
      )

      await user.click(screen.getByRole('button', { name: 'show-details' }))
      expect(screen.queryByRole('button', { name: 'show-details' })).not.toBeInTheDocument()
      expect(screen.getByRole('button', { name: 'hide-details' })).toBeInTheDocument()

      screen.queryAllByLabelText(/song-part-details-/).forEach((d) => expect(d).toBeVisible())

      await user.click(screen.getByRole('button', { name: 'hide-details' }))
      expect(screen.getByRole('button', { name: 'show-details' })).toBeInTheDocument()
      expect(screen.queryByRole('button', { name: 'hide-details' })).not.toBeInTheDocument()
    })

    it('should open song arrangements modal when clicking on song arrangements button', async () => {
      const user = userEvent.setup()

      reduxRender(
        <SongPartsWidget parts={parts} songId={''} settings={emptySongSettings} />
      )

      await user.click(screen.getByRole('button', { name: /song-arrangements/i }))
      expect(await screen.findByRole('dialog', { name: /song arrangements/i })).toBeInTheDocument()
    })

    it('should open add perfect rehearsal popover when on clicking add perfect rehearsal button and send request', async () => {
      const user = userEvent.setup()

      let capturedRequest: AddPerfectSongRehearsalRequest
      server.use(
        http.post('/songs/perfect-rehearsal', async (req) => {
          capturedRequest = (await req.request.json()) as AddPerfectSongRehearsalRequest
          return HttpResponse.json({ message: 'it worked' })
        })
      )

      const songId = 'some-id'

      reduxRender(
        withToastify(
          <SongPartsWidget
            parts={parts}
            songId={songId}
            settings={emptySongSettings}
            defaultSongArrangementId={'id'}
          />
        )
      )

      await user.click(screen.getByRole('button', { name: 'add-perfect-rehearsal' }))

      expect(await screen.findByRole('dialog')).toBeInTheDocument()
      expect(screen.getByText(/increase parts' rehearsals .* occurrences/i)).toBeInTheDocument()

      await user.click(screen.getByRole('button', { name: 'confirm' }))

      expect(await screen.findByText(/perfect rehearsal added/i)).toBeInTheDocument()
      expect(capturedRequest).toStrictEqual({ id: songId })
    })
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

    reduxRender(
      <SongPartsWidget parts={parts} songId={songId} settings={emptySongSettings} />
    )

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

  it('should show add new song part card and open add new song part, when there are no parts', async () => {
    const user = userEvent.setup()

    reduxRender(<SongPartsWidget parts={[]} songId={''} settings={emptySongSettings} />)

    expect(screen.getByLabelText('add-new-song-part-card')).toBeInTheDocument()
    await user.click(screen.getByLabelText('add-new-song-part-card'))
    expect(screen.getByLabelText('add-new-song-part')).toBeInTheDocument()
  })

  it('should show toast when adding 1 rehearsal to part and dismiss on when adding to another part', async () => {
    const user = userEvent.setup()

    const part1 = parts[0]
    const part2 = parts[1]

    reduxRender(
      withToastify(
        <SongPartsWidget parts={parts} songId={''} settings={emptySongSettings} />
      )
    )

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

    reduxRender(<SongPartsWidget parts={parts} songId={''} settings={emptySongSettings} />)

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
