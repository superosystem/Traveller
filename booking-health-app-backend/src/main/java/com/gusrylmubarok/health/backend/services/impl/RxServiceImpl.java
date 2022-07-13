package com.gusrylmubarok.health.backend.services.impl;

import com.gusrylmubarok.health.backend.domain.Rx;
import com.gusrylmubarok.health.backend.repositories.RxDAO;
import com.gusrylmubarok.health.backend.services.RxService;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
public class RxServiceImpl implements RxService {

    final static Logger logger = LoggerFactory.getLogger(RxServiceImpl.class);

    @Autowired private RxDAO rxDAO;

    @Override
    public List<Rx> findByDoctorId(int id) {
        return rxDAO.findByDoctorId(id);
    }

    @Override
    public void save(Rx rx) {
        rxDAO.save(rx);
    }

    @Override
    public List<Rx> findByPatientId(int id) {
        return rxDAO.findByUserId(id);
    }

}
