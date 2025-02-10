import { emptySongSection, reduxRender, withToastify } from '../../../test-utils.tsx'
import EditSongSectionsOccurrencesModal from './EditSongSectionsOccurrencesModal.tsx'
import { SongSection } from '../../../types/models/Song.ts'
import { screen, within } from '@testing-library/react'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import { userEvent } from '@testing-library/user-event'
import { UpdateSongSectionsOccurrencesRequest } from '../../../types/requests/SongRequests.ts'

describe('Edit Song Sections Occurrences Modal', () => {
  const server = setupServer()

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  const sections: SongSection[] = [
    {
      ...emptySongSection,
      id: '1',
      name: 'Sec 1',
      occurrences: 0,
      songSectionType: {
        id: '',
        name: 'Chorus'
      }
    },
    {
      ...emptySongSection,
      id: '2',
      name: 'Sec 2',
      occurrences: 5,
      songSectionType: {
        id: '',
        name: 'Typ'
      }
    }
  ]

  it('should render', async () => {
    const user = userEvent.setup()

    reduxRender(
      <EditSongSectionsOccurrencesModal
        opened={true}
        onClose={() => {}}
        sections={sections}
        songId={''}
      />
    )

    expect(screen.getByRole('dialog', { name: /edit sections' occurrences/i })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /edit sections' occurrences/i })).toBeInTheDocument()

    sections.forEach((section) => {
      const sectionItem = screen.getByLabelText(`section-${section.name}`)
      expect(sectionItem).toBeInTheDocument()
      expect(within(sectionItem).getByRole('textbox', { name: 'type' })).toBeInTheDocument()
      expect(within(sectionItem).getByRole('textbox', { name: 'type' })).toHaveAttribute(
        'readonly',
        ''
      )
      expect(within(sectionItem).getByRole('textbox', { name: 'type' })).toHaveValue(
        section.songSectionType.name
      )

      expect(within(sectionItem).getByRole('textbox', { name: 'name' })).toBeInTheDocument()
      expect(within(sectionItem).getByRole('textbox', { name: 'name' })).toHaveAttribute(
        'readonly',
        ''
      )
      expect(within(sectionItem).getByRole('textbox', { name: 'name' })).toHaveValue(section.name)

      expect(within(sectionItem).getByRole('textbox', { name: 'occurrences' })).toBeInTheDocument()
      expect(within(sectionItem).getByRole('textbox', { name: 'occurrences' })).toHaveValue(
        section.occurrences.toString()
      )

      expect(
        within(sectionItem).getByRole('button', { name: 'decrease-occurrences' })
      ).toBeInTheDocument()
      if (section.occurrences === 0) {
        expect(
          within(sectionItem).getByRole('button', { name: 'decrease-occurrences' })
        ).toBeDisabled()
      }
      expect(
        within(sectionItem).getByRole('button', { name: 'increase-occurrences' })
      ).toBeInTheDocument()
    })

    expect(screen.getByRole('button', { name: /save/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /save/i })).toHaveAttribute('data-disabled', 'true')
    await user.hover(screen.getByRole('button', { name: /save/i }))
    expect(await screen.findByText(/need to make a change/i)).toBeInTheDocument()
  })

  it('should send request to update the occurrences', async () => {
    const user = userEvent.setup()

    const onClose = vitest.fn()
    const songId = 'some-song-id'

    let capturedRequest: UpdateSongSectionsOccurrencesRequest
    server.use(
      http.put('/songs/sections/occurrences', async (req) => {
        capturedRequest = (await req.request.json()) as UpdateSongSectionsOccurrencesRequest
        return HttpResponse.json({ message: 'OK' })
      })
    )

    reduxRender(
      withToastify(
        <EditSongSectionsOccurrencesModal
          opened={true}
          onClose={onClose}
          sections={sections}
          songId={songId}
        />
      )
    )

    await user.click(
      within(screen.getByLabelText(`section-${sections[0].name}`)).getByRole('button', {
        name: 'increase-occurrences'
      })
    )
    await user.click(screen.getByRole('button', { name: /save/i }))

    expect(onClose).toHaveBeenCalledOnce()
    expect(capturedRequest).toStrictEqual({
      songId: songId,
      sections: [
        {
          id: sections[0].id,
          occurrences: sections[0].occurrences + 1 // sole change
        },
        {
          id: sections[1].id,
          occurrences: sections[1].occurrences
        }
      ]
    })
    expect(screen.getByText(/occurrences updated/i)).toBeInTheDocument()
  })

  it('should disable the save button when no change has been made', async () => {
    const user = userEvent.setup()

    reduxRender(
      <EditSongSectionsOccurrencesModal
        opened={true}
        onClose={() => {}}
        sections={sections}
        songId={''}
      />
    )

    const saveButton = screen.getByRole('button', { name: /save/i })

    const section1 = screen.getByLabelText(`section-${sections[0].name}`)
    const section2 = screen.getByLabelText(`section-${sections[1].name}`)

    await user.click(within(section1).getByRole('button', { name: 'increase-occurrences' }))
    expect(saveButton).not.toHaveAttribute('data-disabled')

    await user.click(within(section2).getByRole('button', { name: 'increase-occurrences' }))
    expect(saveButton).not.toHaveAttribute('data-disabled')

    await user.click(within(section1).getByRole('button', { name: 'decrease-occurrences' }))
    expect(saveButton).not.toHaveAttribute('data-disabled')

    await user.click(within(section2).getByRole('button', { name: 'decrease-occurrences' }))
    expect(saveButton).toHaveAttribute('data-disabled', 'true')
  })
})
