package com.gusrylmubarok.springbasic.data.cyclic;

import lombok.extern.slf4j.Slf4j;

import javax.annotation.PostConstruct;
import javax.annotation.PreDestroy;

/**
 * @author Gusryl Mubarok
 * @email gusrylmubarok@gmail.com
 */

@Slf4j
public class Server {

    @PostConstruct
    public void start(){
        log.info("Start Server");
    }

    @PreDestroy
    public void stop(){
        log.info("Stop Server");
    }
}
