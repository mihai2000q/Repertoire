import { emptySongSection, reduxRender, withToastify } from '../../test-utils.tsx'
import SongSectionsCard from './SongSectionsCard.tsx'
import { SongSection } from '../../types/models/Song.ts'
import { fireEvent, screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { expect } from 'vitest'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import {
  AddPerfectSongRehearsalRequest,
  MoveSongSectionRequest
} from '../../types/requests/SongRequests.ts'

describe('Song Sections Card', () => {
  const sections: SongSection[] = [
    {
      ...emptySongSection,
      id: '1',
      name: 'Chorus 1',
      rehearsals: 0,
      confidence: 0,
      progress: 0,
      songSectionType: {
        id: '',
        name: 'Chorus'
      }
    },
    {
      ...emptySongSection,
      id: '2',
      name: 'James Solo',
      rehearsals: 7,
      confidence: 50,
      progress: 163,
      songSectionType: {
        id: '',
        name: 'Solo'
      }
    },
    {
      ...emptySongSection,
      id: '3',
      name: 'James Riff',
      rehearsals: 1,
      confidence: 36,
      progress: 40,
      songSectionType: {
        id: '',
        name: 'Riff'
      }
    }
  ]

  const handlers = [
    http.get('/songs/sections/types', async () => {
      return HttpResponse.json([])
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', () => {
    reduxRender(<SongSectionsCard sections={sections} songId={''} />)

    expect(screen.getByText(/sections/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'add-new-section' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'show-details' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'edit-occurrences' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'edit-occurrences' })).not.toBeDisabled()
    expect(screen.getByRole('button', { name: 'add-perfect-rehearsal' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'add-perfect-rehearsal' })).not.toBeDisabled()

    const renderedSections = screen.getAllByLabelText(/song-section-(?!details)/)
    for (let i = 0; i < sections.length; i++) {
      expect(renderedSections[i]).toHaveAccessibleName(`song-section-${sections[i].name}`)
    }
    screen.queryAllByLabelText(/song-section-details-/).forEach((d) => expect(d).not.toBeVisible())
  })

  it('should disable a few options when there are no sections', () => {
    reduxRender(<SongSectionsCard sections={[]} songId={''} />)

    expect(screen.getByRole('button', { name: 'edit-occurrences' })).toBeDisabled()
    expect(screen.getByRole('button', { name: 'add-perfect-rehearsal' })).toBeDisabled()
  })

  describe('on toolbar options', () => {
    it('should open add new song section when clicking on add new section button', async () => {
      const user = userEvent.setup()

      reduxRender(<SongSectionsCard sections={sections} songId={''} />)

      await user.click(screen.getByRole('button', { name: 'add-new-section' }))
      expect(screen.getByLabelText('add-new-song-section')).toBeInTheDocument()
    })

    it('should show details when clicking on show details', async () => {
      const user = userEvent.setup()

      reduxRender(<SongSectionsCard sections={sections} songId={''} />)

      await user.click(screen.getByRole('button', { name: 'show-details' }))
      expect(screen.queryByRole('button', { name: 'show-details' })).not.toBeInTheDocument()
      expect(screen.getByRole('button', { name: 'hide-details' })).toBeInTheDocument()

      screen.queryAllByLabelText(/song-section-details-/).forEach((d) => expect(d).toBeVisible())

      await user.click(screen.getByRole('button', { name: 'hide-details' }))
      expect(screen.getByRole('button', { name: 'show-details' })).toBeInTheDocument()
      expect(screen.queryByRole('button', { name: 'hide-details' })).not.toBeInTheDocument()
    })

    it("should open edit song sections' occurrences when clicking on edit sections' occurrences button", async () => {
      const user = userEvent.setup()

      reduxRender(<SongSectionsCard sections={sections} songId={''} />)

      await user.click(screen.getByRole('button', { name: 'edit-occurrences' }))
      expect(
        await screen.findByRole('dialog', { name: /edit sections' occurrences/i })
      ).toBeInTheDocument()
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

      reduxRender(withToastify(<SongSectionsCard sections={sections} songId={songId} />))

      await user.click(screen.getByRole('button', { name: 'add-perfect-rehearsal' }))

      expect(await screen.findByRole('dialog')).toBeInTheDocument()
      expect(screen.getByText(/increase sections' rehearsals/i)).toBeInTheDocument()
      expect(screen.getByRole('button', { name: 'cancel-perfect-rehearsal' })).toBeInTheDocument()
      expect(screen.getByRole('button', { name: 'confirm-perfect-rehearsal' })).toBeInTheDocument()

      await user.click(screen.getByRole('button', { name: 'confirm-perfect-rehearsal' }))

      expect(await screen.findByText(/perfect rehearsal added/i)).toBeInTheDocument()
      expect(capturedRequest).toStrictEqual({ id: songId })
    })
  })

  it.skip('should be able to reorder sections', async () => {
    const section = sections[0]
    const overSection = sections[2]
    const songId = 'some-id'

    let capturedRequest: MoveSongSectionRequest
    server.use(
      http.put('/songs/sections/types', async (req) => {
        capturedRequest = (await req.request.json()) as MoveSongSectionRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    reduxRender(<SongSectionsCard sections={sections} songId={songId} />)

    fireEvent.mouseDown(screen.getByLabelText(`song-section-${section.name}`))
    fireEvent.dragStart(screen.getByLabelText(`song-section-${section.name}`))
    fireEvent.dragOver(screen.getByLabelText(`song-section-${overSection.name}`))
    fireEvent.drop(screen.getByLabelText(`song-section-${overSection.name}`))
    fireEvent.mouseUp(screen.getByLabelText(`song-section-${overSection.name}`))

    const renderedSections = await screen.findAllByLabelText(/song-section-(?!details)/)
    const expectedSections = [sections[1], sections[2], sections[0]]

    for (let i = 0; i < sections.length; i++) {
      expect(renderedSections[i]).toHaveAccessibleName(`song-section-${expectedSections[i].name}`)
    }

    expect(capturedRequest).toStrictEqual({
      id: section.id,
      overId: overSection.id,
      songId: songId
    })
  })

  it('should show add new song section card and open add new song section, when there are no sections', async () => {
    const user = userEvent.setup()

    reduxRender(<SongSectionsCard sections={[]} songId={''} />)

    expect(screen.getByLabelText('add-new-song-section-card')).toBeInTheDocument()
    await user.click(screen.getByLabelText('add-new-song-section-card'))
    expect(screen.getByLabelText('add-new-song-section')).toBeInTheDocument()
  })
})
