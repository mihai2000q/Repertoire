import { SongArrangement, SongSectionOccurrences } from '../../../types/models/Song.ts'
import {
  emptySongArrangement,
  emptySongSection,
  reduxRender,
  withToastify
} from '../../../test-utils.tsx'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import SongArrangementsModal from './SongArrangementsModal.tsx'
import { screen, waitFor, within } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import {
  BulkUpdateSongArrangementsRequest,
  UpdateDefaultSongArrangementRequest
} from '../../../types/requests/SongRequests.ts'

describe('Song Arrangements Modal', () => {
  const arrangements: SongArrangement[] = [
    {
      ...emptySongArrangement,
      id: '1',
      name: 'Partial Rehearsal',
      sectionOccurrences: [
        {
          section: {
            ...emptySongSection,
            id: '1',
            name: 'The Verse',
            songSectionType: { id: '1', name: 'Verse' }
          },
          occurrences: 1
        },
        {
          section: {
            ...emptySongSection,
            id: '2',
            name: 'The Chorus',
            songSectionType: { id: '2', name: 'Chorus' }
          },
          occurrences: 2
        },
        {
          section: {
            ...emptySongSection,
            id: '3',
            name: 'The Solo',
            songSectionType: { id: '3', name: 'Solo' }
          },
          occurrences: 0
        }
      ]
    },
    {
      ...emptySongArrangement,
      id: '2',
      name: 'Perfect Rehearsal',
      sectionOccurrences: [
        {
          section: {
            ...emptySongSection,
            id: '1',
            name: 'The Verse',
            songSectionType: { id: '1', name: 'Verse' }
          },
          occurrences: 2
        },
        {
          section: {
            ...emptySongSection,
            id: '2',
            name: 'The Chorus',
            songSectionType: { id: '2', name: 'Chorus' }
          },
          occurrences: 3
        },
        {
          section: {
            ...emptySongSection,
            id: '3',
            name: 'The Solo',
            songSectionType: { id: '3', name: 'Solo' }
          },
          occurrences: 1
        }
      ]
    }
  ]

  const selectedArrangement = arrangements[0]

  const handlers = [
    http.get('/songs/arrangements', async () => {
      return HttpResponse.json(arrangements)
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', async () => {
    const user = userEvent.setup()

    reduxRender(<SongArrangementsModal opened={true} onClose={vi.fn()} songId={''} />)

    // header title
    expect(screen.getByTestId('song-arrangements-loader')).toBeInTheDocument()
    expect(await screen.findByRole('dialog', { name: /song arrangements/i })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /song arrangements/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'add-new-arrangement' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'set-default' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: selectedArrangement.name })).toBeInTheDocument()

    // sub-header
    expect(screen.getByRole('textbox', { name: 'name' })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: 'name' })).toHaveValue(selectedArrangement.name)
    expect(screen.getByRole('button', { name: 'reset' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'reset' })).toBeDisabled()
    expect(screen.getByRole('button', { name: 'delete' })).toBeInTheDocument()

    // content
    selectedArrangement.sectionOccurrences.forEach((sectionOccurrence) => {
      const sectionItem = screen.getByLabelText(`section-${sectionOccurrence.section.name}`)
      expect(sectionItem).toBeInTheDocument()

      // type
      expect(within(sectionItem).getByRole('textbox', { name: 'section-type' })).toBeInTheDocument()
      expect(within(sectionItem).getByRole('textbox', { name: 'section-type' })).toHaveAttribute(
        'readonly',
        ''
      )
      expect(within(sectionItem).getByRole('textbox', { name: 'section-type' })).toHaveValue(
        sectionOccurrence.section.songSectionType.name
      )

      // name
      expect(within(sectionItem).getByRole('textbox', { name: 'section-name' })).toBeInTheDocument()
      expect(within(sectionItem).getByRole('textbox', { name: 'section-name' })).toHaveAttribute(
        'readonly',
        ''
      )
      expect(within(sectionItem).getByRole('textbox', { name: 'section-name' })).toHaveValue(
        sectionOccurrence.section.name
      )

      // occurrences
      expect(
        within(sectionItem).getByRole('textbox', { name: 'section-occurrences' })
      ).toBeInTheDocument()
      expect(within(sectionItem).getByRole('textbox', { name: 'section-occurrences' })).toHaveValue(
        sectionOccurrence.occurrences.toString()
      )

      // occurrences controls
      expect(
        within(sectionItem).getByRole('button', { name: 'decrease-section-occurrences' })
      ).toBeInTheDocument()
      if (sectionOccurrence.occurrences === 0) {
        expect(
          within(sectionItem).getByRole('button', { name: 'decrease-section-occurrences' })
        ).toBeDisabled()
      }
      expect(
        within(sectionItem).getByRole('button', { name: 'increase-section-occurrences' })
      ).toBeInTheDocument()
    })

    // bottom
    expect(screen.getByRole('button', { name: /save/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /save/i })).toHaveAttribute('data-disabled', 'true')
    await user.hover(screen.getByRole('button', { name: /save/i }))
    expect(await screen.findByText(/need to make a change/i)).toBeInTheDocument()
  })

  it('should render even when there are no arrangements', async () => {
    server.use(
      http.get('/songs/arrangements', async () => {
        return HttpResponse.json([])
      })
    )

    reduxRender(<SongArrangementsModal opened={true} onClose={vi.fn()} songId={''} />)

    expect(screen.getByTestId('song-arrangements-loader')).toBeInTheDocument()
    expect(await screen.findByRole('dialog', { name: /song arrangements/i })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /song arrangements/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'add-new-arrangement' })).toBeInTheDocument()
    expect(await screen.findByText(/there are no arrangements/i)).toBeInTheDocument()
  })

  it('should render even when there are no sections on arrangement', async () => {
    server.use(
      http.get('/songs/arrangements', async () => {
        return HttpResponse.json([{ ...selectedArrangement, sectionOccurrences: [] }])
      })
    )

    reduxRender(<SongArrangementsModal opened={true} onClose={vi.fn()} songId={''} />)

    // normal render as expected but without the sections
    expect(await screen.findByText(/there are no sections/i)).toBeInTheDocument()
  })

  it('should be able to set an arrangement as default', async () => {
    const user = userEvent.setup()

    let capturedRequest: UpdateDefaultSongArrangementRequest
    server.use(
      http.put('/songs/arrangements/default', async (req) => {
        capturedRequest = (await req.request.json()) as UpdateDefaultSongArrangementRequest
        return HttpResponse.json('it worked!')
      })
    )

    const songId = 'songId'

    reduxRender(<SongArrangementsModal opened={true} onClose={vi.fn()} songId={songId} />)

    await user.click(await screen.findByRole('button', { name: 'set-default' }))
    await waitFor(() =>
      expect(capturedRequest).toStrictEqual({ id: selectedArrangement.id, songId: songId })
    )
  })

  it('should be able to unset an arrangement from being default', async () => {
    const user = userEvent.setup()

    let capturedRequest: UpdateDefaultSongArrangementRequest
    server.use(
      http.put('/songs/arrangements/default', async (req) => {
        capturedRequest = (await req.request.json()) as UpdateDefaultSongArrangementRequest
        return HttpResponse.json('it worked!')
      })
    )

    const songId = 'songId'

    reduxRender(
      <SongArrangementsModal
        opened={true}
        onClose={vi.fn()}
        songId={songId}
        defaultId={selectedArrangement.id}
      />
    )

    await user.click(await screen.findByRole('button', { name: 'unset-default' }))
    await waitFor(() => expect(capturedRequest).toStrictEqual({ id: null, songId: songId }))
  })

  it('should keep the arrangement name in sync with the menu', async () => {
    const user = userEvent.setup()

    const newName = 'Arr'

    reduxRender(<SongArrangementsModal opened={true} onClose={vi.fn()} songId={''} />)

    const nameTextBox = await screen.findByRole('textbox', { name: 'name' })

    await user.clear(nameTextBox)
    await user.type(nameTextBox, newName)

    expect(screen.getByRole('button', { name: newName }))
    await user.click(screen.getByRole('button', { name: newName }))
    expect(await screen.findByRole('menu', { name: newName })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: newName })).toBeInTheDocument()
  })

  it('should reset the values to initial state when clicking on the reset button', async () => {
    const user = userEvent.setup()

    reduxRender(<SongArrangementsModal opened={true} onClose={vi.fn()} songId={''} />)

    const nameTextBox = await screen.findByRole('textbox', { name: 'name' })
    const section = screen.getByLabelText(
      `section-${selectedArrangement.sectionOccurrences[0].section.name}`
    )

    await user.type(nameTextBox, 's')
    await user.click(within(section).getByRole('button', { name: 'increase-section-occurrences' }))

    expect(nameTextBox).toHaveValue(selectedArrangement.name + 's')
    expect(within(section).getByRole('textbox', { name: 'section-occurrences' })).toHaveValue(
      ((selectedArrangement.sectionOccurrences[0].occurrences as number) + 1).toString()
    )

    await user.click(screen.getByRole('button', { name: 'reset' }))

    expect(nameTextBox).toHaveValue(selectedArrangement.name)
    expect(within(section).getByRole('textbox', { name: 'section-occurrences' })).toHaveValue(
      selectedArrangement.sectionOccurrences[0].occurrences.toString()
    )
  })

  it('should delete the arrangement and select another one when clicking on the delete button', async () => {
    const user = userEvent.setup()

    const songId = 'songId'

    server.use(
      http.delete(`/songs/arrangements/${selectedArrangement.id}/from/${songId}`, () => {
        return HttpResponse.json('it worked!')
      })
    )

    const newArrangement = arrangements[1]

    reduxRender(
      withToastify(<SongArrangementsModal opened={true} onClose={vi.fn()} songId={songId} />)
    )

    // delete
    await user.click(await screen.findByRole('button', { name: 'delete' }))

    // delete - warning modal
    expect(await screen.findByRole('dialog', { name: /delete.* arrangement/i })).toBeInTheDocument()
    await user.click(screen.getByRole('button', { name: /yes/i }))

    expect(await screen.findByText(/song arrangement deleted/i)).toBeInTheDocument()

    // check re-selection
    expect(screen.getByRole('textbox', { name: 'name' })).toHaveValue(newArrangement.name)
    expect(screen.getByRole('button', { name: newArrangement.name })).toBeInTheDocument()
  })

  it('should disable the save button and reset, when no change has been made', async () => {
    const user = userEvent.setup()

    reduxRender(<SongArrangementsModal opened={true} onClose={vi.fn()} songId={''} />)

    const nameTextBox = await screen.findByRole('textbox', { name: 'name' })
    const saveButton = screen.getByRole('button', { name: /save/i })
    const resetButton = screen.getByRole('button', { name: /reset/i })

    const section = screen.getByLabelText(
      `section-${selectedArrangement.sectionOccurrences[0].section.name}`
    )

    // arrangement 1
    await user.type(nameTextBox, 's')
    expect(saveButton).not.toHaveAttribute('data-disabled')
    expect(resetButton).not.toHaveAttribute('data-disabled')
    await user.type(nameTextBox, '{backspace}')
    expect(saveButton).toHaveAttribute('data-disabled', 'true')
    expect(resetButton).toHaveAttribute('data-disabled', 'true')

    await user.click(within(section).getByRole('button', { name: 'increase-section-occurrences' }))
    expect(saveButton).not.toHaveAttribute('data-disabled')
    expect(resetButton).not.toHaveAttribute('data-disabled')
    await user.click(within(section).getByRole('button', { name: 'decrease-section-occurrences' }))
    expect(saveButton).toHaveAttribute('data-disabled', 'true')
    expect(resetButton).toHaveAttribute('data-disabled', 'true')

    // switch
    await user.click(screen.getByRole('button', { name: selectedArrangement.name }))
    await user.click(screen.getByRole('menuitem', { name: arrangements[1].name }))

    // arrangement 2
    await user.type(nameTextBox, 's')
    expect(saveButton).not.toHaveAttribute('data-disabled')
    expect(resetButton).not.toHaveAttribute('data-disabled')
    await user.type(nameTextBox, '{backspace}')
    expect(saveButton).toHaveAttribute('data-disabled', 'true')
    expect(resetButton).toHaveAttribute('data-disabled', 'true')

    await user.click(within(section).getByRole('button', { name: 'increase-section-occurrences' }))
    expect(saveButton).not.toHaveAttribute('data-disabled')
    expect(resetButton).not.toHaveAttribute('data-disabled')
  })

  it('should re-enable the save button when changes are made, even if the selected arrangement changes', async () => {
    const user = userEvent.setup()

    const newName = 'Arr'

    reduxRender(<SongArrangementsModal opened={true} onClose={vi.fn()} songId={''} />)

    const nameTextBox = await screen.findByRole('textbox', { name: 'name' })
    const saveButton = screen.getByRole('button', { name: /save/i })
    const resetButton = screen.getByRole('button', { name: /reset/i })

    const section = screen.getByLabelText(
      `section-${selectedArrangement.sectionOccurrences[0].section.name}`
    )

    await user.clear(nameTextBox)
    await user.type(nameTextBox, newName)
    await user.click(within(section).getByRole('button', { name: 'increase-section-occurrences' }))
    expect(saveButton).not.toHaveAttribute('data-disabled')
    expect(resetButton).not.toHaveAttribute('data-disabled')

    // switch
    await user.click(screen.getByRole('button', { name: newName }))
    await user.click(await screen.findByRole('menuitem', { name: arrangements[1].name }))

    expect(saveButton).not.toHaveAttribute('data-disabled')
    expect(resetButton).toHaveAttribute('data-disabled', 'true') // hence, they are not related
  })

  it('should send update request with the selected arrangement when making changes and saving them', async () => {
    const user = userEvent.setup()

    const songId = 'songId'

    let capturedRequest: BulkUpdateSongArrangementsRequest
    server.use(
      http.put('/songs/arrangements/bulk', async (req) => {
        capturedRequest = (await req.request.json()) as BulkUpdateSongArrangementsRequest
        return HttpResponse.json('it worked!')
      })
    )

    const newName = 'Arr'
    const newSectionOccurrences = [
      {
        ...selectedArrangement.sectionOccurrences[0],
        occurrences: 5
      },
      {
        ...selectedArrangement.sectionOccurrences[1],
        occurrences: 1
      },
      {
        ...selectedArrangement.sectionOccurrences[2],
        occurrences: 1
      }
    ]

    reduxRender(
      withToastify(<SongArrangementsModal opened={true} onClose={vi.fn()} songId={songId} />)
    )

    const nameTextBox = await screen.findByRole('textbox', { name: 'name' })

    await user.clear(nameTextBox)
    await user.type(nameTextBox, newName)

    for (const sectionOccurrence of newSectionOccurrences) {
      const section = screen.getByLabelText(`section-${sectionOccurrence.section.name}`)

      await user.clear(within(section).getByRole('textbox', { name: 'section-occurrences' }))
      await user.type(
        within(section).getByRole('textbox', { name: 'section-occurrences' }),
        sectionOccurrence.occurrences.toString()
      )
    }

    await user.click(screen.getByRole('button', { name: /save/i }))
    expect(await screen.findByText(/song arrangement updated/i)).toBeInTheDocument()
    expect(capturedRequest).toStrictEqual({
      songId: songId,
      requests: [
        {
          id: selectedArrangement.id,
          name: newName,
          occurrences: newSectionOccurrences.map((so) => ({
            sectionId: so.section.id,
            occurrences: so.occurrences
          }))
        }
      ]
    })
  })

  it('should send update request 2 arrangements when making changes and saving them', async () => {
    const user = userEvent.setup()

    const songId = 'songId'

    let capturedRequest: BulkUpdateSongArrangementsRequest
    server.use(
      http.put('/songs/arrangements/bulk', async (req) => {
        capturedRequest = (await req.request.json()) as BulkUpdateSongArrangementsRequest
        return HttpResponse.json('it worked!')
      })
    )

    const newName = 'Arr'
    const newSectionOccurrences: SongSectionOccurrences[] = [
      {
        ...arrangements[0].sectionOccurrences[0],
        occurrences: 5
      },
      {
        ...arrangements[0].sectionOccurrences[1],
        occurrences: 1
      },
      {
        ...arrangements[0].sectionOccurrences[2],
        occurrences: 1
      }
    ]

    const newName2 = 'Arr2'
    const newSectionOccurrences2: SongSectionOccurrences[] = [
      {
        ...arrangements[1].sectionOccurrences[0],
        occurrences: 5
      },
      {
        ...arrangements[1].sectionOccurrences[1],
        occurrences: 1
      },
      {
        ...arrangements[1].sectionOccurrences[2],
        occurrences: 1
      }
    ]

    reduxRender(
      withToastify(<SongArrangementsModal opened={true} onClose={vi.fn()} songId={songId} />)
    )

    const nameTextBox = await screen.findByRole('textbox', { name: 'name' })

    async function makeChangesOnArrangement(
      name: string,
      sectionOccurrences: SongSectionOccurrences[]
    ) {
      await user.clear(nameTextBox)
      await user.type(nameTextBox, name)

      for (const sectionOccurrence of sectionOccurrences) {
        const section = screen.getByLabelText(`section-${sectionOccurrence.section.name}`)

        await user.clear(within(section).getByRole('textbox', { name: 'section-occurrences' }))
        await user.type(
          within(section).getByRole('textbox', { name: 'section-occurrences' }),
          sectionOccurrence.occurrences.toString()
        )
      }
    }

    // make changes
    await makeChangesOnArrangement(newName, newSectionOccurrences)
    await user.click(screen.getByRole('button', { name: newName }))
    await user.click(await screen.findByRole('menuitem', { name: arrangements[1].name }))
    await makeChangesOnArrangement(newName2, newSectionOccurrences2)

    // save
    await user.click(screen.getByRole('button', { name: /save/i }))

    // check
    expect(await screen.findByText(/song arrangements updated/i)).toBeInTheDocument()
    expect(capturedRequest).toStrictEqual({
      songId: songId,
      requests: [
        {
          id: arrangements[0].id,
          name: newName,
          occurrences: newSectionOccurrences.map((so) => ({
            sectionId: so.section.id,
            occurrences: so.occurrences
          }))
        },
        {
          id: arrangements[1].id,
          name: newName2,
          occurrences: newSectionOccurrences.map((so) => ({
            sectionId: so.section.id,
            occurrences: so.occurrences
          }))
        }
      ]
    })
  })
})
