import { emptySongSection, reduxRender, withToastify } from '../../../test-utils.tsx'
import { Instrument, SongSection, SongSectionType } from '../../../types/models/Song.ts'
import { setupServer } from 'msw/node'
import { fireEvent, screen, waitFor } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { http, HttpResponse } from 'msw'
import { UpdateSongSectionRequest } from '../../../types/requests/SongRequests.ts'
import EditSongSectionModal from './EditSongSectionModal.tsx'
import { BandMember } from '../../../types/models/Artist.ts'

describe('Edit Song Description Modal', () => {
  const bandMembers: BandMember[] = [
    {
      id: '1',
      name: 'Nick',
      roles: [{ id: '1', name: 'Guitarist' }]
    },
    {
      id: '2',
      name: 'Joe',
      roles: [{ id: '2', name: 'Vocalist' }]
    }
  ]

  const instruments: Instrument[] = [
    {
      id: '1',
      name: 'Guitar'
    },
    {
      id: '2',
      name: 'Piano'
    }
  ]

  const sectionTypes: SongSectionType[] = [
    {
      id: '1',
      name: 'Solo'
    },
    {
      id: '2',
      name: 'Riff'
    }
  ]

  const section: SongSection = {
    ...emptySongSection,
    id: 'some-id',
    name: 'section 1',
    songSectionType: sectionTypes[1],
    rehearsals: 12,
    confidence: 50
  }

  const handlers = [
    http.get(`/songs/sections/types`, () => {
      return HttpResponse.json(sectionTypes)
    }),
    http.get(`/songs/instruments`, () => {
      return HttpResponse.json(instruments)
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', async () => {
    const user = userEvent.setup()

    reduxRender(
      <EditSongSectionModal
        opened={true}
        onClose={() => {}}
        section={section}
        bandMembers={bandMembers}
      />
    )

    expect(screen.getByRole('dialog', { name: /edit song section/i })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /edit song section/i })).toBeInTheDocument()

    expect(screen.getByRole('textbox', { name: /name/i })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: /name/i })).not.toBeInvalid()
    expect(screen.getByRole('textbox', { name: /name/i })).toHaveValue(section.name)

    expect(screen.getByRole('textbox', { name: /type/i })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: /type/i })).not.toBeInvalid()
    expect(await screen.findByRole('textbox', { name: /type/i })).toHaveValue(
      section.songSectionType.name
    )

    expect(screen.getByRole('textbox', { name: /rehearsals/i })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: /rehearsals/i })).not.toBeInvalid()
    expect(screen.getByRole('textbox', { name: /rehearsals/i })).toHaveValue(
      section.rehearsals.toString()
    )

    expect(screen.getByRole('textbox', { name: /band member/i })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: /band member/i })).toHaveValue(
      section.bandMember?.name ?? ''
    )

    expect(screen.getByRole('textbox', { name: /instrument/i })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: /instrument/i })).toHaveValue(
      section.instrument?.name ?? ''
    )

    expect(screen.getByRole('slider', { name: /confidence/i })).toBeInTheDocument()
    expect(screen.getByRole('slider', { name: /confidence/i })).toHaveValue(section.confidence)

    expect(screen.getByRole('button', { name: /save/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /save/i })).toHaveAttribute('data-disabled', 'true')
    await user.hover(screen.getByRole('button', { name: /save/i }))
    expect(await screen.findByText(/need to make a change/i)).toBeInTheDocument()
  })

  it('should send update request when the field values have changed', async () => {
    const user = userEvent.setup()

    const newName = 'New Section Name'
    const newType = sectionTypes[0]
    const newRehearsals = 23
    const newConfidence = 82
    const onClose = vitest.fn()

    let capturedRequest: UpdateSongSectionRequest
    server.use(
      http.put('/songs/sections', async (req) => {
        capturedRequest = (await req.request.json()) as UpdateSongSectionRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    const [{ rerender }] = reduxRender(
      withToastify(
        <EditSongSectionModal
          opened={true}
          onClose={onClose}
          section={section}
          bandMembers={bandMembers}
        />
      )
    )

    const nameField = screen.getByRole('textbox', { name: /name/i })
    const typeField = screen.getByRole('textbox', { name: /type/i })
    const rehearsalsField = screen.getByRole('textbox', { name: /rehearsals/i })
    const confidenceField = screen.getByRole('slider', { name: /confidence/i })
    const saveButton = screen.getByRole('button', { name: /save/i })

    await user.clear(nameField)
    await user.type(nameField, newName)

    await user.click(typeField)
    await user.click(await screen.findByText(newType.name))

    await user.clear(rehearsalsField)
    await user.type(rehearsalsField, newRehearsals.toString())

    for (let i = section.confidence; i < newConfidence; i++) {
      fireEvent.keyDown(confidenceField, { key: 'ArrowRight' })
    }

    expect(saveButton).not.toHaveAttribute('data-disabled')
    await user.click(screen.getByRole('button', { name: /save/i }))

    expect(capturedRequest).toStrictEqual({
      id: section.id,
      name: newName,
      typeId: newType.id,
      rehearsals: newRehearsals,
      confidence: newConfidence
    })
    expect(onClose).toHaveBeenCalledOnce()

    expect(await screen.findByText(`${newName} updated!`)).toBeInTheDocument()

    rerender(
      <EditSongSectionModal
        opened={true}
        onClose={onClose}
        section={{
          ...section,
          name: newName,
          songSectionType: newType,
          rehearsals: newRehearsals,
          confidence: newConfidence
        }}
        bandMembers={bandMembers}
      />
    )
    expect(screen.getByRole('button', { name: /save/i })).toHaveAttribute('data-disabled', 'true')
  })

  it('should send update request when the band member and instruments changed', async () => {
    const user = userEvent.setup()

    const newBandMember = bandMembers[0]
    const newInstrument = instruments[0]
    const onClose = vitest.fn()

    let capturedRequest: UpdateSongSectionRequest
    server.use(
      http.put('/songs/sections', async (req) => {
        capturedRequest = (await req.request.json()) as UpdateSongSectionRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    const [{ rerender }] = reduxRender(
      withToastify(
        <EditSongSectionModal
          opened={true}
          onClose={onClose}
          section={section}
          bandMembers={bandMembers}
        />
      )
    )

    const bandMemberField = screen.getByRole('textbox', { name: /band member/i })
    const instrumentField = screen.getByRole('textbox', { name: /instrument/i })
    const saveButton = screen.getByRole('button', { name: /save/i })

    await user.click(bandMemberField)
    await user.click(await screen.findByText(newBandMember.name))

    await user.click(instrumentField)
    await user.click(await screen.findByText(newInstrument.name))

    expect(saveButton).not.toHaveAttribute('data-disabled')
    await user.click(screen.getByRole('button', { name: /save/i }))

    expect(capturedRequest).toStrictEqual({
      id: section.id,
      name: section.name,
      typeId: section.songSectionType.id,
      rehearsals: section.rehearsals,
      confidence: section.confidence,
      bandMemberId: newBandMember.id,
      instrumentId: newInstrument.id
    })
    expect(onClose).toHaveBeenCalledOnce()

    expect(await screen.findByText(`${section.name} updated!`)).toBeInTheDocument()

    rerender(
      <EditSongSectionModal
        opened={true}
        onClose={onClose}
        section={{
          ...section,
          bandMember: newBandMember,
          instrument: newInstrument
        }}
        bandMembers={bandMembers}
      />
    )
    expect(screen.getByRole('button', { name: /save/i })).toHaveAttribute('data-disabled', 'true')
  })

  it('should keep the save button disabled when the field values have not changed', async () => {
    const user = userEvent.setup()

    reduxRender(
      <EditSongSectionModal
        opened={true}
        onClose={() => {}}
        section={section}
        bandMembers={bandMembers}
      />
    )

    const nameField = screen.getByRole('textbox', { name: /name/i })
    const typeField = screen.getByRole('textbox', { name: /type/i })
    const rehearsalsField = screen.getByRole('textbox', { name: /rehearsals/i })
    const bandMemberField = screen.getByRole('textbox', { name: /band member/i })
    const instrumentField = screen.getByRole('textbox', { name: /instrument/i })
    const confidenceField = screen.getByRole('slider', { name: /confidence/i })
    const saveButton = screen.getByRole('button', { name: /save/i })

    // change name
    await user.clear(nameField)
    await user.type(nameField, section.name + '1')
    expect(saveButton).not.toHaveAttribute('data-disabled')

    // reset name
    await user.clear(nameField)
    await user.type(nameField, section.name)
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    // change type
    await user.click(typeField)
    await user.click(await screen.findByText(sectionTypes[0].name))
    expect(saveButton).not.toHaveAttribute('data-disabled')

    // reset type
    await user.click(typeField)
    await user.click(await screen.findByText(section.songSectionType.name))
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    // change rehearsals
    await user.clear(rehearsalsField)
    await user.type(rehearsalsField, section.rehearsals.toString() + '1')
    expect(saveButton).not.toHaveAttribute('data-disabled')

    // reset rehearsals
    await user.clear(rehearsalsField)
    await user.type(rehearsalsField, section.rehearsals.toString())
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    // change band member
    await user.click(bandMemberField)
    await user.click(await screen.findByText(bandMembers[0].name))
    expect(saveButton).not.toHaveAttribute('data-disabled')

    // reset band member
    await user.click(bandMemberField)
    await user.click(await screen.findByText(section.bandMember?.name ?? bandMembers[0].name))
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    // change an instrument
    await user.click(instrumentField)
    await user.click(await screen.findByText(instruments[0].name))
    expect(saveButton).not.toHaveAttribute('data-disabled')

    // reset instrument
    await user.click(instrumentField)
    await user.click(await screen.findByText(section.instrument?.name ?? instruments[0].name))
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    // change confidence
    fireEvent.keyDown(confidenceField, { key: 'ArrowRight' })
    expect(saveButton).not.toHaveAttribute('data-disabled')

    // reset confidence
    fireEvent.keyDown(confidenceField, { key: 'ArrowLeft' })
    expect(saveButton).toHaveAttribute('data-disabled', 'true')
  })

  it('should validate fields', async () => {
    const user = userEvent.setup()

    reduxRender(
      <EditSongSectionModal
        opened={true}
        onClose={() => {}}
        section={section}
        bandMembers={bandMembers}
      />
    )

    const nameField = screen.getByRole('textbox', { name: /name/i })
    const rehearsalsField = screen.getByRole('textbox', { name: /rehearsals/i })

    // invalidate name
    await user.clear(nameField)
    expect(nameField).toBeInvalid()

    // invalidate rehearsals - cannot be empty
    await user.clear(rehearsalsField)
    expect(rehearsalsField).toBeInvalid()

    // reset rehearsals
    await user.type(rehearsalsField, section.rehearsals.toString())
    expect(rehearsalsField).not.toBeInvalid()

    // invalidate rehearsals - cannot be lower than initial value
    await user.clear(rehearsalsField)
    await user.type(rehearsalsField, (section.rehearsals - 1).toString())
    expect(rehearsalsField).toBeInvalid()
  })

  it('should keep fields updated', async () => {
    const [{ rerender }] = reduxRender(
      <EditSongSectionModal
        opened={true}
        onClose={() => {}}
        section={section}
        bandMembers={bandMembers}
      />
    )

    expect(screen.getByRole('textbox', { name: /rehearsals/i })).toHaveValue(
      section.rehearsals.toString()
    )
    expect(screen.getByRole('textbox', { name: /band member/i })).toHaveValue(section.bandMember?.name ?? '')
    expect(screen.getByRole('textbox', { name: /instrument/i })).toHaveValue(section.instrument?.name ?? '')

    const newSection = {
      ...section,
      rehearsals: section.rehearsals + 1,
      instrument: instruments[0],
      bandMember: bandMembers[1]
    }

    rerender(
      <EditSongSectionModal
        opened={true}
        onClose={() => {}}
        section={newSection}
        bandMembers={bandMembers}
      />
    )

    expect(screen.getByRole('textbox', { name: /rehearsals/i })).toHaveValue(
      newSection.rehearsals.toString()
    )
    await waitFor(() =>
      expect(screen.getByRole('textbox', { name: /instrument/i })).toHaveValue(
        newSection.instrument.name
      )
    )
    await waitFor(() =>
      expect(screen.getByRole('textbox', { name: /band member/i })).toHaveValue(
        newSection.bandMember.name
      )
    )
  })
})
