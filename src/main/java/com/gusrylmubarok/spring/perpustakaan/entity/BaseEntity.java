package com.gusrylmubarok.spring.perpustakaan.entity;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;
import org.hibernate.annotations.DynamicUpdate;

import javax.persistence.*;
import java.io.Serializable;
import java.util.Date;

@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@MappedSuperclass
@DynamicUpdate
@SuppressWarnings("unchecked")
public abstract class BaseEntity<T>  implements Serializable {

    private static final long serialVersionUID = 2836987502693397892L;

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Column(name = "created_time")
    private Date createdTime;

    @Column(name = "updated_time")
    private Date updatedTime;

    @PrePersist
    protected void onCreate() {
        setCreatedTime(new Date());
    }

    @PreUpdate
    protected void onUpdate() {
        setUpdatedTime(new Date());
    }

}
