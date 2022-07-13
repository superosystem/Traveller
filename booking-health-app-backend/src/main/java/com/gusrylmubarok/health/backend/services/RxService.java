package com.gusrylmubarok.health.backend.services;

import com.gusrylmubarok.health.backend.domain.Rx;

import java.util.List;

public interface RxService {
    void save(Rx rx);
    List<Rx> findByDoctorId(int doctorId);
    List<Rx> findByPatientId(int userId);
}
