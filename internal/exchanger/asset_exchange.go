package exchanger

func (ex *Exchanger) listenUpdateMeta() {
	ch := ex.mnt.ListenUpdateMeta()
	for {
		select {
		case <-ex.ctx.Done():
			return
		case updateMeta, ok := <-ch:
			ex.logger.Info("Receive to-update meta from monitor")
			if !ok {
				ex.logger.Warn("Unexpected closed channel while listening on update meta")
				return
			}
			if err := ex.syncer.SendUpdateMeta(updateMeta); err != nil {
				ex.logger.Errorf("Send update meta error: %s", err.Error())
				return
			}
			ex.logger.Info("Update meta event successfully")
		}
	}
}

func (ex *Exchanger) listenMintEvent() {
	ch := ex.mnt.ListenLockEvent()
	for {
		select {
		case <-ex.ctx.Done():
			return
		case lockEvent, ok := <-ch:
			ex.logger.Info("Receive lock event from monitor")
			if !ok {
				ex.logger.Warn("Unexpected closed channel while listening on lock event")
				return
			}
			// do handleMissingEvent
			if int64(lockEvent.GetAppchainIndex())-1 > ex.rAppchainIndex {
				ex.handleMissingLockFromMnt(ex.rAppchainIndex, int64(lockEvent.GetAppchainIndex())-1)
			}
			if err := ex.syncer.SendLockEvent(lockEvent); err != nil {
				ex.logger.Errorf("Send lock event error: %s", err.Error())
				return
			}
			ex.rAppchainIndex++
			ex.logger.Info("Lock event successfully")
		}
	}
}

func (ex *Exchanger) listenBurnEventFromSyncer() {
	// start bxhJsonRpc client
	ex.syncer.JsonrpcClient().Start(ex.aRelayIndex)
	ch := ex.syncer.ListenBurn()
	for {
		select {
		case <-ex.ctx.Done():
			return
		case burnEvent, ok := <-ch:
			if !ok {
				ex.logger.Warn("Unexpected closed channel while listening on interchain burn event")
				return
			}
			// do handleMissingEvent
			if int64(burnEvent.GetRelayIndex())-1 > ex.aRelayIndex {
				ex.handleMissingBurnFromSyncer(ex.aRelayIndex, int64(burnEvent.GetRelayIndex())-1)
			}
			// get mutil signs
			burnEvent.MultiSigns, _ = ex.syncer.GetEVMSigns(burnEvent.TxId)
			if err := ex.exec.SendBurnEvent(burnEvent); err != nil {
				// handle sending error
				ex.logger.Errorf("Send unlock event error: %s", err.Error())
				return
			}
			ex.aRelayIndex++
			ex.logger.Info("unlock event successfully")

		}
	}
}