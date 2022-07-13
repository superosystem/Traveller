package com.gusrylmubarok.health.backend.repositories;

import com.gusrylmubarok.health.backend.domain.Rx;

import java.util.List;

public interface RxDAO {
    List<Rx> findByDoctorId(int doctorId);
    List<Rx> findByUserId(int userId);
    Rx save(Rx rx);
}
