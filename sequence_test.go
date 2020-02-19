/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package mirbft

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	pb "github.com/IBM/mirbft/mirbftpb"
)

var _ = Describe("sequence", func() {
	var (
		s *sequence
	)

	BeforeEach(func() {
		s = &sequence{
			myConfig: &Config{
				ID: 1,
			},
			networkConfig: &pb.NetworkConfig{
				Nodes: []uint64{0, 1, 2, 3},
				F:     1,
			},
			epoch:    4,
			seqNo:    5,
			prepares: map[string]map[NodeID]struct{}{},
			commits:  map[string]map[NodeID]struct{}{},
		}
	})

	Describe("allocate", func() {
		It("transitions from Unknown to Allocated", func() {
			actions := s.allocate(
				[]*request{
					{
						//TODO, []byte("msg1"),
					},
					{
						//TODO, []byte("msg2"),
					},
				},
			)

			Expect(actions).To(Equal(&Actions{
				Process: []*Batch{
					{
						SeqNo: 5,
						Epoch: 4,
						// TODO, broken
						// Proposals: [][]byte{
						// []byte("msg1"),
						// []byte("msg2"),
						// },
					},
				},
			}))

			Expect(s.state).To(Equal(Allocated))
			Expect(s.batch).To(Equal(
				[][]byte{
					[]byte("msg1"),
					[]byte("msg2"),
				},
			))
		})

		When("the current state is not Unknown", func() {
			BeforeEach(func() {
				s.state = Prepared
			})

			It("does not transition and instead panics", func() {
				badTransition := func() {
					s.allocate(
						[]*request{
							// TODO, broken
							// []byte("msg1"),
							// []byte("msg2"),
						},
					)
				}
				Expect(badTransition).To(Panic())
				Expect(s.state).To(Equal(Prepared))
			})
		})
	})

	Describe("applyProcessResult", func() {
		BeforeEach(func() {
			s.state = Allocated
			s.batch = []*request{
				// TODO, broken
				// []byte("msg1"),
				// []byte("msg2"),
			}
		})

		It("transitions from Allocated to Preprepared", func() {
			actions := s.applyProcessResult([]byte("digest"), true)
			Expect(actions).To(Equal(&Actions{
				Broadcast: []*pb.Msg{
					{
						Type: &pb.Msg_Prepare{
							Prepare: &pb.Prepare{
								SeqNo:  5,
								Epoch:  4,
								Digest: []byte("digest"),
							},
						},
					},
				},
				QEntries: []*pb.QEntry{
					{
						SeqNo:  5,
						Epoch:  4,
						Digest: []byte("digest"),
						// TODO, broken
						// Proposals: [][]byte{
						// []byte("msg1"),
						// []byte("msg2"),
						// },
					},
				},
			}))
			Expect(s.digest).To(Equal([]byte("digest")))
			Expect(s.state).To(Equal(Preprepared))
			Expect(s.qEntry).To(Equal(&pb.QEntry{
				SeqNo:  5,
				Epoch:  4,
				Digest: []byte("digest"),
				// TODO, broken
				// Proposals: [][]byte{
				// []byte("msg1"),
				// []byte("msg2"),
				// },
			}))

		})

		When("the state is not Allocated", func() {
			BeforeEach(func() {
				s.state = Prepared
			})

			It("does not transition the state and panics", func() {
				badTransition := func() {
					s.applyProcessResult([]byte("digest"), true)
				}
				Expect(badTransition).To(Panic())
				Expect(s.state).To(Equal(Prepared))
			})
		})

		When("when the validation is not successful", func() {
			It("transitions the state to InvalidBatch", func() {
				actions := s.applyProcessResult([]byte("digest"), false)
				Expect(actions).To(Equal(&Actions{}))
				Expect(s.state).To(Equal(Invalid))
				Expect(s.digest).To(Equal([]byte("digest")))
				Expect(s.qEntry).To(Equal(&pb.QEntry{
					SeqNo:  5,
					Epoch:  4,
					Digest: []byte("digest"),
					// TODO, broken
					// Proposals: [][]byte{
					// []byte("msg1"),
					// []byte("msg2"),
					// },
				}))
			})
		})
	})

	Describe("applyPrepareMsg", func() {
		BeforeEach(func() {
			s.state = Preprepared
			s.digest = []byte("digest")
			s.prepares["digest"] = map[NodeID]struct{}{
				1: {},
				2: {},
			}
		})

		It("transitions from Preprepared to Prepared", func() {
			actions := s.applyPrepareMsg(0, []byte("digest"))
			Expect(actions).To(Equal(&Actions{
				Broadcast: []*pb.Msg{
					{
						Type: &pb.Msg_Commit{
							Commit: &pb.Commit{
								SeqNo:  5,
								Epoch:  4,
								Digest: []byte("digest"),
							},
						},
					},
				},
				PEntries: []*pb.PEntry{
					{
						SeqNo:  5,
						Epoch:  4,
						Digest: []byte("digest"),
					},
				},
			}))
		})
	})
})
