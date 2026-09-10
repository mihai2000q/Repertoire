import {
  emptySongArrangement,
  emptySongPart,
  emptySongSettings,
  reduxRender,
  withToastify
} from '../../../../../../test-utils.tsx'
import { SongArrangement, SongPart } from '../../../../../../types/models/Song.ts'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { expect } from 'vitest'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import { AddPerfectSongRehearsalRequest } from '../../../../../../types/requests/SongRequests.ts'
import { createRef } from 'react'
import SongOutlineWidget from './SongOutlineWidget.tsx'

// Mock Main Context
vi.mock('../../../../../../context/MainContext.tsx', () => ({
  useMain: vi.fn(() => ({
    ref: createRef(),
    mainScroll: { ref: createRef() }
  }))
}))

describe('Song Outline Widget', () => {
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

  const arrangements: SongArrangement[] = [{ ...emptySongArrangement, id: '1' }]

  const handlers = [
    http.get('/songs/arrangements', () => {
      return HttpResponse.json(arrangements)
    }),
    http.get('/songs/instruments', () => {
      return HttpResponse.json([])
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
      <SongOutlineWidget
        parts={parts}
        songId={''}
        settings={emptySongSettings}
        defaultSongArrangementId={'1'}
      />
    )

    expect(screen.getByText(/outline/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'add-new-part' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'show-details' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'show-details' })).not.toBeDisabled()
    expect(screen.getByRole('button', { name: 'manage-song-arrangements' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'add-custom-rehearsal' })).toBeInTheDocument()
    expect(await screen.findByRole('button', { name: 'add-custom-rehearsal' })).not.toBeDisabled()
    expect(screen.getByRole('button', { name: 'add-perfect-rehearsal' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'add-perfect-rehearsal' })).not.toBeDisabled()
    expect(screen.getByRole('button', { name: 'settings' })).toBeInTheDocument()
    expect(screen.getByRole('radio', { name: 'sections-view' })).toBeInTheDocument()
    expect(screen.getByRole('radio', { name: 'parts-view' })).toBeInTheDocument()
    expect(screen.getByLabelText('song-parts')).toBeInTheDocument()
  })

  it('should disable a few options when there are no parts', () => {
    reduxRender(<SongOutlineWidget parts={[]} songId={''} settings={emptySongSettings} />)

    expect(screen.getByRole('button', { name: 'show-details' })).toBeDisabled()
    expect(screen.getByRole('button', { name: 'add-custom-rehearsal' })).toBeDisabled()
    expect(screen.getByRole('button', { name: 'add-perfect-rehearsal' })).toBeDisabled()
  })

  it('should disable perfect rehearsal when a default arrangement is not set', () => {
    reduxRender(<SongOutlineWidget parts={parts} songId={''} settings={emptySongSettings} />)

    expect(screen.getByRole('button', { name: 'add-perfect-rehearsal' })).toBeDisabled()
  })

  describe('on toolbar options', () => {
    it('should open add new song part when clicking on add new part button', async () => {
      const user = userEvent.setup()

      reduxRender(<SongOutlineWidget parts={parts} songId={''} settings={emptySongSettings} />)

      await user.click(screen.getByRole('button', { name: 'add-new-part' }))
      expect(screen.getByLabelText('add-new-song-part')).toBeInTheDocument()
    })

    it('should show details when clicking on show details', async () => {
      const user = userEvent.setup()

      reduxRender(<SongOutlineWidget parts={parts} songId={''} settings={emptySongSettings} />)

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

      reduxRender(<SongOutlineWidget parts={parts} songId={''} settings={emptySongSettings} />)

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
          <SongOutlineWidget
            parts={parts}
            songId={songId}
            settings={emptySongSettings}
            defaultSongArrangementId={'1'}
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

  it('should show add new song part card and open add new song part, when there are no parts', async () => {
    const user = userEvent.setup()

    reduxRender(<SongOutlineWidget parts={[]} songId={''} settings={emptySongSettings} />)

    expect(screen.getByLabelText('add-new-song-part-card')).toBeInTheDocument()
    await user.click(screen.getByLabelText('add-new-song-part-card'))
    expect(screen.getByLabelText('add-new-song-part')).toBeInTheDocument()
  })
})
