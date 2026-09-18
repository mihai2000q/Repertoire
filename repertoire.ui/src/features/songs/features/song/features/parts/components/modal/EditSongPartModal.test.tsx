import { emptySongPart, reduxRender, withToastify } from '../../../../../../../../test-utils.tsx'
import { Instrument, SongPart } from '../../../../../../../../types/models/Song.ts'
import { setupServer } from 'msw/node'
import { fireEvent, screen, waitFor } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { http, HttpResponse } from 'msw'
import { UpdateSongPartRequest } from '../../types/requests/SongPartRequests.ts'
import EditSongPartModal from './EditSongPartModal.tsx'
import { BandMember } from '../../../../../../../../types/models/Artist.ts'
import { ReactNode } from 'react'

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

  const part: SongPart = {
    ...emptySongPart,
    id: 'some-id',
    name: 'part 1',
    rehearsals: 12,
    confidence: 50,
    bandMembers: [bandMembers[1]]
  }

  const handlers = [
    http.get(`/songs/instruments`, () => {
      return HttpResponse.json(instruments)
    })
  ]

  const server = setupServer(...handlers)

  const render = (ui: ReactNode) =>
    reduxRender(ui, { song: { songId: '', isArtistBand: true, artistBandMembers: bandMembers } })

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', async () => {
    const user = userEvent.setup()

    const [{ rerender }] = render(
      <EditSongPartModal opened={true} onClose={() => {}} part={part} sectionId={'section-1'} />
    )

    expect(screen.getByRole('dialog', { name: /edit song part/i })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /edit song part/i })).toBeInTheDocument()

    expect(screen.getByRole('textbox', { name: /name/i })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: /name/i })).not.toBeInvalid()
    expect(screen.getByRole('textbox', { name: /name/i })).toHaveValue(part.name)

    expect(screen.getByRole('textbox', { name: /rehearsals/i })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: /rehearsals/i })).not.toBeInvalid()
    expect(screen.getByRole('textbox', { name: /rehearsals/i })).toHaveValue(
      part.rehearsals.toString()
    )

    expect(screen.getByRole('combobox', { name: /band member/i })).toBeInTheDocument()
    expect(screen.getByRole('combobox', { name: /band member/i })).toHaveValue(
      part.bandMembers[0]?.name ?? ''
    )

    expect(screen.getByRole('combobox', { name: /instrument/i })).toBeInTheDocument()
    expect(screen.getByRole('combobox', { name: /instrument/i })).toHaveValue(
      part.instrument?.name ?? ''
    )

    expect(screen.getByRole('slider', { name: /confidence/i })).toBeInTheDocument()
    expect(screen.getByRole('slider', { name: /confidence/i })).toHaveValue(part.confidence)

    expect(screen.getByRole('button', { name: /save/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /save/i })).toBeDisabled()
    await user.hover(screen.getByRole('button', { name: /save/i }))
    expect(await screen.findByText(/need to make a change/i)).toBeInTheDocument()

    // without section
    rerender(
      <EditSongPartModal opened={true} onClose={() => {}} part={part} />
    )

    expect(screen.queryByRole('combobox', { name: /band member/i })).not.toBeInTheDocument()
  })

  it('should send update request when the field values have changed', async () => {
    const user = userEvent.setup()

    const newName = 'New Part Name'
    const newRehearsals = 23
    const newConfidence = 82
    const onClose = vitest.fn()

    let capturedRequest: UpdateSongPartRequest
    server.use(
      http.put('/songs/parts', async (req) => {
        capturedRequest = (await req.request.json()) as UpdateSongPartRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    const [{ rerender }] = render(
      withToastify(
        <EditSongPartModal opened={true} onClose={onClose} part={part} />
      )
    )

    const nameField = screen.getByRole('textbox', { name: /name/i })
    const rehearsalsField = screen.getByRole('textbox', { name: /rehearsals/i })
    const confidenceField = screen.getByRole('slider', { name: /confidence/i })
    const saveButton = screen.getByRole('button', { name: /save/i })

    await user.clear(nameField)
    await user.type(nameField, newName)

    await user.clear(rehearsalsField)
    await user.type(rehearsalsField, newRehearsals.toString())

    for (let i = part.confidence; i < newConfidence; i++) {
      fireEvent.keyDown(confidenceField, { key: 'ArrowRight' })
    }

    expect(saveButton).not.toBeDisabled()
    await user.click(screen.getByRole('button', { name: /save/i }))

    expect(capturedRequest).toStrictEqual({
      id: part.id,
      name: newName,
      rehearsals: newRehearsals,
      confidence: newConfidence
    })
    expect(onClose).toHaveBeenCalledOnce()

    expect(await screen.findByText(`${newName} updated!`)).toBeInTheDocument()

    rerender(
      <EditSongPartModal
        opened={true}
        onClose={onClose}
        part={{
          ...part,
          name: newName,
          rehearsals: newRehearsals,
          confidence: newConfidence
        }}
      />
    )
    expect(screen.getByRole('button', { name: /save/i })).toBeDisabled()
  })

  it('should send update request when the band member and instruments changed', async () => {
    const user = userEvent.setup()

    const sectionId = 'section-1'
    const newBandMember = bandMembers[0]
    const newInstrument = instruments[0]
    const onClose = vitest.fn()

    let capturedRequest: UpdateSongPartRequest
    server.use(
      http.put('/songs/parts', async (req) => {
        capturedRequest = (await req.request.json()) as UpdateSongPartRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    const [{ rerender }] = render(
      withToastify(
        <EditSongPartModal opened={true} onClose={onClose} part={part} sectionId={sectionId} />
      )
    )

    const bandMemberField = screen.getByRole('combobox', { name: /band member/i })
    const instrumentField = screen.getByRole('combobox', { name: /instrument/i })
    const saveButton = screen.getByRole('button', { name: /save/i })

    await user.click(bandMemberField)
    await user.clear(bandMemberField)
    await user.click(await screen.findByText(newBandMember.name))

    await user.click(instrumentField)
    await user.click(await screen.findByText(newInstrument.name))

    expect(saveButton).not.toBeDisabled()
    await user.click(screen.getByRole('button', { name: /save/i }))

    expect(capturedRequest).toStrictEqual({
      id: part.id,
      name: part.name,
      rehearsals: part.rehearsals,
      confidence: part.confidence,
      bandMemberId: newBandMember.id,
      instrumentId: newInstrument.id,
      sectionId: sectionId
    })
    expect(onClose).toHaveBeenCalledOnce()

    expect(await screen.findByText(`${part.name} updated!`)).toBeInTheDocument()

    rerender(
      <EditSongPartModal
        opened={true}
        onClose={onClose}
        part={{
          ...part,
          bandMembers: [newBandMember],
          instrument: newInstrument
        }}
        sectionId={'section-1'}
      />
    )
    expect(screen.getByRole('button', { name: /save/i })).toBeDisabled()
  })

  it('should keep the save button disabled when the field values have not changed', async () => {
    const user = userEvent.setup()

    render(
      <EditSongPartModal opened={true} onClose={() => {}} part={part} sectionId={'section-1'} />
    )

    const nameField = screen.getByRole('textbox', { name: /name/i })
    const rehearsalsField = screen.getByRole('textbox', { name: /rehearsals/i })
    const bandMemberField = screen.getByRole('combobox', { name: /band member/i })
    const instrumentField = screen.getByRole('combobox', { name: /instrument/i })
    const confidenceField = screen.getByRole('slider', { name: /confidence/i })
    const saveButton = screen.getByRole('button', { name: /save/i })

    // change name
    await user.clear(nameField)
    await user.type(nameField, part.name + '1')
    expect(saveButton).not.toBeDisabled()

    // reset name
    await user.clear(nameField)
    await user.type(nameField, part.name)
    expect(saveButton).toBeDisabled()

    // change rehearsals
    await user.clear(rehearsalsField)
    await user.type(rehearsalsField, part.rehearsals.toString() + '1')
    expect(saveButton).not.toBeDisabled()

    // reset rehearsals
    await user.clear(rehearsalsField)
    await user.type(rehearsalsField, part.rehearsals.toString())
    expect(saveButton).toBeDisabled()

    // change band member
    await user.click(bandMemberField)
    await user.clear(bandMemberField)
    await user.click(await screen.findByText(bandMembers[0].name))
    expect(saveButton).not.toBeDisabled()

    // reset band member
    await user.clear(bandMemberField)
    await user.click(await screen.findByText(part.bandMembers[0].name))
    expect(saveButton).toBeDisabled()

    // change an instrument
    await user.click(instrumentField)
    await user.click(await screen.findByText(instruments[0].name))
    expect(saveButton).not.toBeDisabled()

    // reset instrument
    await user.click(instrumentField)
    await user.click(await screen.findByText(part.instrument?.name ?? instruments[0].name))
    expect(saveButton).toBeDisabled()

    // change confidence
    fireEvent.keyDown(confidenceField, { key: 'ArrowRight' })
    expect(saveButton).not.toBeDisabled()

    // reset confidence
    fireEvent.keyDown(confidenceField, { key: 'ArrowLeft' })
    expect(saveButton).toBeDisabled()
  })

  it('should validate fields', async () => {
    const user = userEvent.setup()

    render(
      <EditSongPartModal opened={true} onClose={() => {}} part={part} />
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
    await user.type(rehearsalsField, part.rehearsals.toString())
    expect(rehearsalsField).not.toBeInvalid()

    // invalidate rehearsals - cannot be lower than initial value
    await user.clear(rehearsalsField)
    await user.type(rehearsalsField, (part.rehearsals - 1).toString())
    expect(rehearsalsField).toBeInvalid()
  })

  it('should keep fields updated', async () => {
    const [{ rerender }] = render(
      <EditSongPartModal opened={true} onClose={() => {}} part={part} sectionId={'section-1'} />
    )

    expect(screen.getByRole('textbox', { name: /rehearsals/i })).toHaveValue(
      part.rehearsals.toString()
    )
    expect(screen.getByRole('combobox', { name: /band member/i })).toHaveValue(
      part.bandMembers[0]?.name ?? ''
    )
    expect(screen.getByRole('combobox', { name: /instrument/i })).toHaveValue(
      part.instrument?.name ?? ''
    )

    const newPart = {
      ...part,
      rehearsals: part.rehearsals + 1,
      instrument: instruments[0],
      bandMembers: [bandMembers[1]]
    }

    rerender(
      <EditSongPartModal opened={true} onClose={() => {}} part={newPart} sectionId={'section-1'} />
    )

    expect(screen.getByRole('textbox', { name: /rehearsals/i })).toHaveValue(
      newPart.rehearsals.toString()
    )
    await waitFor(() =>
      expect(screen.getByRole('combobox', { name: /instrument/i })).toHaveValue(
        newPart.instrument.name
      )
    )
    await waitFor(() =>
      expect(screen.getByRole('combobox', { name: /band member/i })).toHaveValue(
        newPart.bandMembers[0].name
      )
    )
  })
})
